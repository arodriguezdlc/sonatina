package terraformcli

import (
	"os/exec"

	"github.com/sirupsen/logrus"
)

func (t *Terraform) Init(path string) error {
	args := []string{}
	args = append(args, "init")
	args = append(args, t.initDefaultOptions().array()...)
	logrus.WithField("args", args).Info("executing terraform command")
	cmd := exec.Command(t.BinaryPath(), args...)
	cmd.Dir = path
	return t.runPrintingAll(cmd)
}

func (t *Terraform) initDefaultOptions() *options {
	return &options{
		option{
			key:   "backend",
			value: "false",
		},
		option{
			key:   "input",
			value: "false",
		},
		option{
			key:   "no-color",
			value: "",
		},
	}
}
