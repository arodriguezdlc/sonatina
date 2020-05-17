package terraformcli

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

type command struct {
}

func (c *command) runPrintingAll(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "error executing terraform init")
	}

	return nil
}
