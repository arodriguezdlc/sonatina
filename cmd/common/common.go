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
