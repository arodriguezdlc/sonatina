package deployment

import (
	"path/filepath"

	"github.com/arodriguezdlc/sonatina/gitw"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// CTD represents a Code Tree Definition
type CTD struct {
	fs   afero.Fs
	path string
	git  *gitw.Command

	Name     string
	RepoURL  string
	RepoPath string

	main    *main
	modules *modules
	vtd     *VTD
}

type main struct {
	fs   afero.Fs
	path string
}

type modules struct {
	fs   afero.Fs
	path string
}

// NewCTD returns an initialized CTD struct
func NewCTD(fs afero.Fs, path string, name string, repoURL string, repoPath string) *CTD {
	git, _ := gitw.NewCommand(fs, path)

	ctd := &CTD{
		fs:   fs,
		path: path,
		git:  git,

		Name:     name,
		RepoURL:  repoURL,
		RepoPath: repoPath,

		main: &main{
			fs:   fs,
			path: filepath.Join(path, "main"),
		},
		modules: &modules{
			fs:   fs,
			path: filepath.Join(path, "modules"),
		},
		vtd: NewVTD(fs, filepath.Join(path, "vtd")),
	}
	return ctd
}

// ListMainGlobalFiles returns all TF files from the global main folder
func (ctd *CTD) ListMainGlobalFiles() ([]string, error) {
	slice, err := afero.Glob(ctd.main.fs, filepath.Join(ctd.main.globalPath(), "/*.tf"))
	if err != nil {
		return slice, errors.Wrap(err, "cannot list main global files")
	}
	return slice, nil
}

// ListMainUserFiles returns all TF files from the user main folder
func (ctd *CTD) ListMainUserFiles() ([]string, error) {
	slice, err := afero.Glob(ctd.main.fs, filepath.Join(ctd.main.userPath(), "/*.tf"))
	if err != nil {
		return slice, errors.Wrap(err, "cannot list main user files")
	}
	return slice, nil
}

// ListModules returns all modules defined on a CTD
func (ctd *CTD) ListModules() ([]string, error) {
	slice, err := afero.Glob(ctd.modules.fs, filepath.Join(ctd.modules.path, "/*"))
	if err != nil {
		return slice, errors.Wrap(err, "cannot list modules")
	}
	return slice, nil
}

// Clone executes a `git clone` command equivalent to get the CTD repository
func (ctd *CTD) Clone() error {
	return ctd.git.Clone(ctd.RepoURL)
}

// Pull executes a `git pull` command equivalent to update the CTD repository
// with the last changes
func (ctd *CTD) Pull() error {
	ctd.git.Pull("origin", "master") // TODO: use a specified branch
	return nil
}

// Checkout executes a `git checkout` command equivalent to point the CTD
// repository to a specified commit, branch or tag
func (ctd *CTD) Checkout() error {
	// TODO: support for checkout to a specified version, branch or commit
	return nil
}

func (m *main) globalPath() string {
	return filepath.Join(m.path, "global")
}

func (m *main) userPath() string {
	return filepath.Join(m.path, "user")
}
