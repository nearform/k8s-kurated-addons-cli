package docker

import (
	"fmt"
	"github.com/nearform/k8s-kurated-addons-cli/src/utils/defaults"
)

type DockerImage struct {
	Registry  string
	Name      string
	Directory string
	Tag       string
	Port      string
}

func (dockerImage DockerImage) RemoteTag() string {
	if dockerImage.Directory != defaults.ProjectDirectory {
		return fmt.Sprintf("%s/%s/%s:%s", dockerImage.Registry, dockerImage.Name, dockerImage.Directory, dockerImage.Tag)
	}
	return fmt.Sprintf("%s/%s:%s", dockerImage.Registry, dockerImage.Name, dockerImage.Tag)
}

func (dockerImage DockerImage) LocalTag() string {
	if dockerImage.Directory != defaults.ProjectDirectory {
		return fmt.Sprintf("%s/%s:%s", dockerImage.Name, dockerImage.Directory, dockerImage.Tag)
	}
	return fmt.Sprintf("%s:%s", dockerImage.Name, dockerImage.Tag)
}
