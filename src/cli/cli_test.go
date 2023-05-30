package cli

import (
	"os"
	"testing"
)

var root = "../../"

func TestEnvConfig(t *testing.T) {
    cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error %s", err)
	}

	os.Args = []string{"./bin/kka-cli", "template"}

    cli := CLI{
        CWD: cwd,
    }

    err = cli.Run()

    if (err != nil) {

    }
}