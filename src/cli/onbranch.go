package cli

import (
	"fmt"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/nearform/k8s-kurated-addons-cli/src/utils"
	"github.com/urfave/cli/v2"
)

func (c CLI) buildPushDeploy(cCtx *cli.Context) error {
	err := c.Build(cCtx)
	if err != nil {
		return fmt.Errorf("building %v", err)
	}
	if cCtx.Bool(stopOnBuildFlag) {
		return err
	}

	err = c.Push(cCtx)
	if err != nil {
		return fmt.Errorf("pushing %v", err)
	}
	if cCtx.Bool(stopOnPushFlag) {
		return err
	}
	return c.Deploy(cCtx)
}

func (c CLI) OnBranchCMD() *cli.Command {
	flags := []cli.Flag{}
	flags = append(flags, c.CommandFlags(Kubernetes)...)
	flags = append(flags, c.CommandFlags(Build)...)
	flags = append(flags, c.CommandFlags(Registry)...)
	flags = append(flags, []cli.Flag{
		&cli.BoolFlag{
			Name:  stopOnBuildFlag,
			Value: false,
		},
		&cli.BoolFlag{
			Name:  stopOnPushFlag,
			Value: false,
		},
		&cli.BoolFlag{
			Name:  "clean",
			Value: false,
		},
	}...)
	return &cli.Command{
		Name:  "onbranch",
		Usage: "deploy the application as a knative service",
		Flags: flags,
		Action: func(cCtx *cli.Context) error {
			if cCtx.Bool("clean") {
				return c.Delete(cCtx)
			}

			return c.buildPushDeploy(cCtx)
		},
		Before: func(ctx *cli.Context) error {
			if err := c.loadFlagsFromConfig(ctx); err != nil {
				return err
			}

			wd, err := os.Getwd()

			if err != nil {
				return err
			}

			repo, err := git.PlainOpen(wd)
			if err != nil {
				return err
			}

			head, err := repo.Head()
			if err != nil {
				return err
			}

			branchName := head.Name().Short()
			c.Logger.Infof("Using branch %v as version and namespace", branchName)

			ctx.Set(appVersionFlag, utils.EncodeRFC1123(branchName))
			ctx.Set(namespaceFlag, utils.EncodeRFC1123(branchName))

			ignoredFlags := []string{}
			if ctx.Bool(stopOnBuildFlag) {
				ignoredFlags = append(ignoredFlags, []string{registryPasswordFlag, registryUserFlag}...)
			}
			if ctx.Bool(stopOnPushFlag) {
				ignoredFlags = append(ignoredFlags, []string{endpointFlag, tokenFlag, caCRTFlag, namespaceFlag}...)
			}

			return c.checkRequiredFlags(ctx, ignoredFlags)
		},
	}
}
