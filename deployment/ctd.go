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
// TODO
func NewCTD(fs afero.Fs, path string, repoURL string, repoPath string) *CTD {
	git, _ := gitw.NewCommand(fs, path)

	ctd := &CTD{
		fs:   fs,
		path: path,
		git:  git,

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
func (ctd *CTD) ListMainUserFiles(user string) ([]string, error) {
	slice, err := afero.Glob(ctd.main.fs, filepath.Join(ctd.main.userPath(user), "/*.tf"))
	if err != nil {
		return slice, errors.Wrap(err, "cannot list main user files")
	}
	return slice, nil
}

// ListModules returns all modules defined on a CTD
func (ctd *CTD) ListModules() ([]string, error) {
	// TODO: modules could be obtained from other repos
	slice, err := afero.Glob(ctd.modules.fs, filepath.Join(ctd.modules.path, "/*"))
	if err != nil {
		return slice, errors.Wrap(err, "cannot list modules")
	}
	return slice, nil
}

func (ctd *CTD) Clone() error {
	return ctd.git.Clone(ctd.RepoURL)
}

func (ctd *CTD) Pull() error {
	// TODO
	return nil
}

func (ctd *CTD) Checkout() error {
	// TODO
	return nil
}

func (m *main) globalPath() string {
	return filepath.Join(m.path, "global")
}

func (m *main) userPath(user string) string {
	return filepath.Join(m.path, "user", user)
}
