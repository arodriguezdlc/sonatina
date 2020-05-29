package terraformcli

import (
	"os/exec"

	"github.com/sirupsen/logrus"
)

func (t *Terraform) Apply(path string, varFiles []string, stateFile string) error {
	args := []string{}
	args = append(args, "apply")
	args = append(args, t.applyDefaultOptions().array()...)
	args = append(args, t.varFilesOptions(varFiles).array()...)
	args = append(args, t.stateFileOption(stateFile).render())
	logrus.WithField("args", args).Info("executing terraform command")

	cmd := exec.Command(t.BinaryPath(), args...)
	cmd.Dir = path

	return t.runPrintingAll(cmd)
}

func (t *Terraform) applyDefaultOptions() *options {
	return &options{
		option{
			key:   "auto-approve",
			value: "",
		},
		option{
			key:   "input",
			value: "false",
		},
		option{
			key:   "no-color",
			value: "",
		},
		option{
			key:   "compact-warnings",
			value: "",
		},
	}
}
