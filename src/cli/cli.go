package cli

import (
	"embed"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/nearform/k8s-kurated-addons-cli/src/services/project"
	"k8s.io/utils/strings/slices"

	"github.com/charmbracelet/log"
	"github.com/nearform/k8s-kurated-addons-cli/src/services/docker"
	"github.com/nearform/k8s-kurated-addons-cli/src/utils/logger"
	"github.com/urfave/cli/v2"
)

type CLI struct {
	Resources     embed.FS
	CWD           string
	DockerService docker.DockerService
	Logger        *log.Logger
	project       project.Project
	dockerImage   docker.DockerImage
	Writer        io.Writer
}

func (c CLI) baseBeforeFunc(ctx *cli.Context) error {
	if err := c.loadFlagsFromConfig(ctx); err != nil {
		return err
	}

	if err := c.checkRequiredFlags(ctx, []string{}); err != nil {
		return err
	}
	return nil
}

func (c *CLI) init(cCtx *cli.Context) {
	appName := cCtx.String(appNameFlag)
	version := cCtx.String(appVersionFlag)
	projectDirectory := cCtx.String(projectDirectoryFlag)
	absProjectDirectory, err := filepath.Abs(cCtx.String(projectDirectoryFlag))

	if err != nil {
		c.Logger.Warnf("could not get abs of %s", projectDirectory)
		absProjectDirectory = projectDirectory
	}

	project := project.New(
		appName,
		projectDirectory,
		cCtx.String(runtimeVersionFlag),
		version,
		c.Resources,
	)

	dockerImageName := appName
	invalidBases := []string{".", string(os.PathSeparator)}
	base := filepath.Base(absProjectDirectory)
	if !slices.Contains(invalidBases, base) && base != dockerImageName {
		dockerImageName = appName + "/" + base
	}

	dockerImage := docker.DockerImage{
		Registry:  cCtx.String(repoNameFlag),
		Name:      dockerImageName,
		Directory: absProjectDirectory,
		Tag:       version,
	}

	dockerService, err := docker.New(project, dockerImage, cCtx.String(dockerFileNameFlag))
	if err != nil {
		logger.PrintError("Error creating docker service", err)
	}

	c.DockerService = dockerService
	c.dockerImage = dockerImage
	c.project = project
}

func (c *CLI) getProject(cCtx *cli.Context) *project.Project {
	if (c.project == project.Project{}) {
		c.init(cCtx)
	}
	return &c.project
}

func (c CLI) Run(args []string) error {
	app := &cli.App{
		Name:  "k8s kurated addons",
		Usage: "kka-cli",
		Flags: c.CommandFlags(App),
		Commands: []*cli.Command{
			c.BuildCMD(),
			c.PushCMD(),
			c.DeployCMD(),
			c.DeleteCMD(),
			c.OnMainCMD(),
			c.OnBranchCMD(),
			c.TemplateCMD(),
			c.InitCMD(),
		},
		Before: func(ctx *cli.Context) error {
			if err := c.loadFlagsFromConfig(ctx); err != nil {
				return err
			}

			projectDirectory := ctx.String(projectDirectoryFlag)
			absProjectDirectory, err := filepath.Abs(projectDirectory)

			if err != nil {
				return err
			}

			if ctx.String(appNameFlag) == "" {
				ctx.Set(appNameFlag, path.Base(absProjectDirectory))
			}

			if err := c.checkRequiredFlags(ctx, []string{}); err != nil {
				return err
			}
			return nil
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	return app.Run(args)
}
