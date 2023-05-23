package cli

import (
	"github.com/docker/docker/api/types"
	"github.com/urfave/cli/v2"
)

func (c *CLI) Push(cCtx *cli.Context) error {
	c.init(cCtx)
	c.DockerService.AuthConfig = types.AuthConfig{
		Username: cCtx.String("registry-user"),
		Password: cCtx.String("registry-password"),
	}
	return c.DockerService.Push()
}

func (c *CLI) PushCMD() *cli.Command {
	return &cli.Command{
		Name:   "push",
		Usage:  "push the container image to a registry",
		Flags:  Flags(Registry),
		Action: c.Push,
	}
}
