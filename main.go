package main

import (
	"embed"
	"log"
	"os"

	"github.com/nearform/k8s-kurated-addons-cli/src/cli"
	"github.com/nearform/k8s-kurated-addons-cli/src/utils/logger"
)

//go:embed assets
var resources embed.FS

func main() {
	logger.PrintInfo("nearForm: k8s kurated addons CLI")

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cli := cli.CLI{
		Resources: resources,
		CWD:       cwd,
	}


	if err := cli.Run(); err != nil {
        log.Fatal(err)
    }
}
