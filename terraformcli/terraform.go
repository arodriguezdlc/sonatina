package terraformcli

import (
	"github.com/spf13/afero"
)

type Terraform struct {
	fs   afero.Fs
	path string

	binary
	command
}

// New constructs a new Terraform struct and returns it
func New(fs afero.Fs, path string, version string, arch string) (*Terraform, error) {
	binary := binary{
		fs:      fs,
		path:    path,
		version: version,
		arch:    arch,
	}
	terraform := &Terraform{
		fs:   fs,
		path: path,

		binary: binary,
	}

	err := fs.MkdirAll(path, 0755)
	if err != nil {
		return terraform, err
	}

	ok, err := terraform.checkBinary()
	if err != nil {
		return terraform, err
	}
	if !ok {
		err = terraform.getBinary()
		if err != nil {
			return terraform, err
		}
	}

	return terraform, nil
}
