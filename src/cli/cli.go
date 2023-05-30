package cli

import (
	"embed"
	"os"
	"path"

	"github.com/nearform/k8s-kurated-addons-cli/src/services/project"

	"github.com/nearform/k8s-kurated-addons-cli/src/services/docker"
	"github.com/nearform/k8s-kurated-addons-cli/src/utils/defaults"
	"github.com/nearform/k8s-kurated-addons-cli/src/utils/logger"

	"github.com/urfave/cli/v2"
    "github.com/joho/godotenv"
)

type CLI struct {
	Resources     embed.FS
	CWD           string
	DockerService docker.DockerService
	Project       project.Project
	dockerImage   docker.DockerImage
}

func (c *CLI) init(cCtx *cli.Context) {

	repoName := cCtx.String("repo-name")
	dockerFileName := cCtx.String("dockerfile-name")
	appName := cCtx.String("app-name")
	version := cCtx.String("app-version")
	projectDirectory := cCtx.String("project-directory")

	project := project.New(
		appName,
		projectDirectory,
		cCtx.String("runtime-version"),
		version,
		c.Resources,
	)

	dockerImage := docker.DockerImage{
		Registry:  repoName,
		Name:      appName,
		Directory: projectDirectory,
		Tag:       version,
	}

	dockerService, err := docker.New(project, dockerImage, dockerFileName)
	if err != nil {
		logger.PrintError("Error creating docker service", err)
	}

	c.DockerService = dockerService
	c.dockerImage = dockerImage
	c.Project = project
}

func (c *CLI) getProject(cCtx *cli.Context) *project.Project {
	if (c.Project == project.Project{}) {
		c.init(cCtx)
	}
	return &c.Project
}

func (c CLI) Run() error {
    configFilePath := os.Getenv("KKA_CONFIG_FILE")
    if (configFilePath != "") {
        err := godotenv.Load(configFilePath)
        if (err != nil) {
            logger.PrintError("Unable to open config file", err)
        }
    }

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "app-name",
			Usage:   "The name of the app",
			Value:   path.Base(c.CWD),
			EnvVars: []string{"KKA_APP_NAME"},
		},
		&cli.StringFlag{
			Name:    "app-version",
			Usage:   "The version of your application",
			Value:   "latest",
			EnvVars: []string{"KKA_VERSION"},
		},
		&cli.StringFlag{
			Name:    "project-directory",
			Usage:   "The directory in which your Dockerfile lives",
			Value:   defaults.ProjectDirectory,
			EnvVars: []string{"KKA_PROJECT_DIRECTORY"},
		},
		&cli.StringFlag{
			Name:    "repo-name",
			Usage:   "The base address of the container repository",
			Value:   defaults.RepoName,
			EnvVars: []string{"KKA_REPO_NAME"},
		},
		&cli.StringFlag{
			Name:    "dockerfile-name",
			Usage:   "The name of the Dockerfile",
			Value:   defaults.DockerfileName,
			EnvVars: []string{"KKA_DOCKERFILE_NAME"},
		},
	}

	app := &cli.App{
		Name:  "k8s kurated addons",
		Usage: "kka-cli",
		Flags: flags,
		Commands: []*cli.Command{
			c.BuildCMD(),
			c.PushCMD(),
			c.DeployCMD(),
			c.DeleteCMD(),
			c.OnMainCMD(),
			c.TemplateCMD(),
			c.InitCMD(),
		},
	}

	return app.Run(os.Args)
}
