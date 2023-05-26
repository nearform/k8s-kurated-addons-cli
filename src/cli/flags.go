package cli

import (
	"github.com/urfave/cli/v2"
)

type FlagsType string

const (
	Build      FlagsType = "build"
	Kubernetes FlagsType = "kubernetes"
	Registry   FlagsType = "registry"
	Init       FlagsType = "init"
)

func Flags(command FlagsType) []cli.Flag {
	flags := map[FlagsType]([]cli.Flag){
		Build: []cli.Flag{
			&cli.StringFlag{
				Name:     "runtime-version",
				EnvVars:  []string{"KKA_RUNTIME_VERSION"},
				Category: "build",
			},
		},
		Kubernetes: []cli.Flag{
			&cli.StringFlag{
				Name:     "endpoint",
				EnvVars:  []string{"KKA_ENDPOINT"},
				Required: true,
				Category: "deploy",
			},
			&cli.StringFlag{
				Name:     "token",
				EnvVars:  []string{"KKA_TOKEN"},
				Required: true,
				Category: "deploy",
			},
			&cli.StringFlag{
				Name:     "ca-crt",
				EnvVars:  []string{"KKA_CA_CERT"},
				Required: true,
				Category: "deploy",
			},
			&cli.StringFlag{
				Name:     "app-port",
				EnvVars:  []string{"KKA_APP_PORT"},
				Required: false,
				Category: "deploy",
			},
		},
		Registry: []cli.Flag{
			&cli.StringFlag{
				Name:     "registry-user",
				EnvVars:  []string{"KKA_REGISTRY_USER"},
				Required: true,
				Category: "registry",
			},
			&cli.StringFlag{
				Name:     "registry-password",
				EnvVars:  []string{"KKA_REGISTRY_PASSWORD"},
				Required: true,
				Category: "registry",
			},
		},
	}

	return flags[command]
}
