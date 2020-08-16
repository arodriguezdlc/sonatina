package common

import (
	"github.com/arodriguezdlc/sonatina/deployment"
	"github.com/arodriguezdlc/sonatina/terraformcli"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

// Fs is the afero.Fs used by sonatina. It's in package common to be shared across
// all commands.
var Fs afero.Fs

// InitializeTerraform initializes terraform object with correct parameters. This initialization
// is used across all commands that needs it.
func InitializeTerraform(deployment deployment.Deployment) (*terraformcli.Terraform, error) {
	terraformPath, err := homedir.Expand(viper.GetString("TerraformPath"))
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't get terraform file path")
	}

	return terraformcli.New(Fs, terraformPath, deployment.TerraformVersion(), viper.GetString("BinaryArch"))
}

// GetCurrentDeployment returns the current deployment name saved on the current file, or
// the override provided as argument if is set.
func GetCurrentDeployment(override string) (string, error) {
	if override != "" {
		return override, nil
	}

	filename, err := homedir.Expand("~/.sonatina/current")
	if err != nil {
		return "", errors.Wrap(err, "couldn't expand directory")
	}

	deployData, err := afero.ReadFile(Fs, filename)
	if err != nil {
		return "", errors.Wrap(err, "couldn't read current deployment")
	}

	return string(deployData), nil
}

// SetCurrentDeployment sets a current deployment saving it to the current file.
func SetCurrentDeployment(deployName string) error {
	filename, err := homedir.Expand("~/.sonatina/current")
	if err != nil {
		return errors.Wrap(err, "couldn't expand directory")
	}

	err = afero.WriteFile(Fs, filename, []byte(deployName), 0644)
	if err != nil {
		return errors.Wrap(err, "couldn't write current deployment")
	}

	return nil
}
