package k8s

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"text/template"
	"time"

	"github.com/nearform/k8s-kurated-addons-cli/src/services/docker"
	"github.com/nearform/k8s-kurated-addons-cli/src/services/project"

	"github.com/nearform/k8s-kurated-addons-cli/src/utils/logger"

	corev1 "k8s.io/api/core/v1"
	apimachineryErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	certutil "k8s.io/client-go/util/cert"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
	servingv1client "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"
)

func Config(endpoint string, token string, caCrt []byte) (*rest.Config, error) {
	if _, err := certutil.NewPoolFromBytes(caCrt); err != nil {
		return nil, fmt.Errorf("Expected to load root CA from bytes, but got err: %v", err)
	}
	return &rest.Config{
		// TODO: switch to using cluster DNS.
		Host: endpoint,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: caCrt,
		},
		BearerToken: string(token),
	}, nil
}

func loadManifest(project *project.Project, dockerImage docker.DockerImage) (*servingv1.Service, error) {
	knativeTemplate := path.Join("assets", "knative", "service.yaml.tmpl")
	template, err := template.ParseFS(project.Resources, knativeTemplate)
	if err != nil {
		return nil, fmt.Errorf("error reading the knative service yaml: %v", err)
	}

	output := &bytes.Buffer{}
	// TODO replace map[string]string{} with proper values
	if err = template.Execute(output, dockerImage); err != nil {
		return nil, err
	}

	data := output.Bytes()

	err = servingv1.AddToScheme(scheme.Scheme)
	if err != nil {
		return nil, fmt.Errorf("error adding Knative Serving scheme: %v", err)
	}

	// Decode the YAML data into a Knative Service object
	obj, _, err := scheme.Codecs.UniversalDeserializer().Decode(data, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error decoding YAML: %v", err)
	}

	service, ok := obj.(*servingv1.Service)
	if !ok {
		return nil, fmt.Errorf("decoded object is not a Knative Service: %v", obj)
	}

	return service, nil
}

func Apply(namespace string, config *rest.Config, project *project.Project, dockerImage docker.DockerImage) error {
	logger.PrintInfo("Deploying Knative service to " + config.Host)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	fmt.Println(dockerImage.RemoteTag())
	service, err := loadManifest(project, dockerImage)
	if err != nil {
		return err
	}

	// Create a new Knative Serving client
	servingClient, err := servingv1client.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("Error creating the knative client %v", err)
	}

	service.ObjectMeta.Namespace = namespace
	service.ObjectMeta.Name = project.Name

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("Creating Kubernetes client %v", err)
	}

	_, err = client.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: service.ObjectMeta.Namespace,
		},
	}, metav1.CreateOptions{})

	if err != nil && !apimachineryErrors.IsAlreadyExists(err) {
		return fmt.Errorf("cannot create namespace %s, failed with %v", service.ObjectMeta.Namespace, err)
	}

	getService, err := servingClient.Services(service.ObjectMeta.Namespace).Get(ctx, service.ObjectMeta.Name, metav1.GetOptions{})
	var deployedService *servingv1.Service
	if err != nil {
		deployedService, err = servingClient.Services(service.ObjectMeta.Namespace).Create(ctx, service, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("Creating Knative service %v", err)
		}
	} else {
		deployedService, err = servingClient.Services(service.ObjectMeta.Namespace).Update(ctx, getService, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("Updating Knative service %v", err)
		}
	}

	fmt.Printf("Knative service %q deployed successfully in namespace %s.\n", deployedService.GetObjectMeta().GetName(), service.ObjectMeta.Namespace)
	return nil
}

func Clean(namespace string, config *rest.Config, project *project.Project) error {
	logger.PrintInfo("Deleting Knative service from " + config.Host)
	ctx := context.Background()

	// Create a new Knative Serving client
	servingClient, err := servingv1client.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("Error creating the knative client %v", err)
	}

	err = servingClient.Services(namespace).Delete(ctx, project.Name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("deleting the service: %v", err)
	}

	logger.PrintInfo("Deleted Knative service " + project.Name)
	return nil
}
