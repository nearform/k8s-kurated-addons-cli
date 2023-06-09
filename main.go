package main

import (
	"embed"
	"os"

	"github.com/charmbracelet/log"
	"github.com/nearform/k8s-kurated-addons-cli/src/cli"
)

//go:embed assets
var resources embed.FS

func main() {
	cli := cli.CLI{
		Resources: resources,
		Logger: log.NewWithOptions(os.Stderr, log.Options{
			Level:           log.ParseLevel(os.Getenv("KKA_LOG_LEVEL")),
			ReportCaller:    true,
			ReportTimestamp: true,
		}),
		Writer: os.Stdout,
	}

	if err := cli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
