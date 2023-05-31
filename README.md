# k8s Kurated Addons Cli

This project is intended for building an executable file to run on various CI/CD platforms in order to get a project deployed to a containerised environment

### Pre-requisites

1. GoLang

    You can install it with your prefered package manager or usign asfd with `asdf install`.

2. Docker (or similar solutions)  

    Here you can find a list of possible candidates:
    - [Docker](https://docs.docker.com/engine/install/) ( cross-platform, paid solution )
    - [Rancher Desktop](https://rancherdesktop.io/) ( cross-platform, FOSS )
    - [lima](https://github.com/lima-vm/lima) + [nerdctl](https://github.com/containerd/nerdctl) ( macOS only )


### Build & run example Nodejs application from code

1. Run `go run main.go --project-directory example build`
2. Run `docker run k8s-kurated-addons-cli/example:latest`
3. You should see the hello world message.
4. Remove the image `docker image rmi -f k8s-kurated-addons-cli/example:latest`

### Build the executable

In order to build the executable you simply need to run 

```bash
make build
```

You will be able to run the executable from 

```bash
./bin/kka-cli
```
