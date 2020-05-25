package deployment

import (
	"path/filepath"

	"github.com/spf13/afero"
)

type VTD struct {
	fs   afero.Fs
	path string

	config  *config
	flavour *flavour
	static  *static
}

type config struct {
	path string
}

type flavour struct {
	path string
}

type static struct {
	path string
}

func NewVTD(fs afero.Fs, path string) *VTD {
	return &VTD{
		fs:   fs,
		path: path,

		config: &config{
			path: filepath.Join(path, "config"),
		},
		flavour: &flavour{
			path: filepath.Join(path, "config"),
		},
		static: &static{
			path: filepath.Join(path, "static"),
		},
	}
}

func (c *config) globalFile() string {
	return filepath.Join(c.path, "global.tfvars")
}

func (c *config) userFile() string {
	return filepath.Join(c.path, "user.tfvars")
}

func (f *flavour) globalFile(flavour string) string {
	return filepath.Join(f.path, "global", flavour+".tfvars")
}

func (f *flavour) userFile(flavour string) string {
	return filepath.Join(f.path, "user", flavour+".tfvars")
}

func (s *static) globalFile() string {
	return filepath.Join(s.path, "global.tfvars")
}

func (s *static) userFile() string {
	return filepath.Join(s.path, "user.tfvars")
}
