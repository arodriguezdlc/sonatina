package deployment

import (
	"path/filepath"

	"github.com/spf13/afero"
)

// CTD represents a Code Tree Definition
type CTD struct {
	fs   afero.Fs
	path string

	RepoURL  string
	RepoPath string

	main    main
	modules modules
	vtd     VTD
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
func NewCTD(fs afero.Fs, path string, repoURL string, repoPath string) (*CTD, error) {
	ctd := &CTD{
		fs:   fs,
		path: path,

		RepoURL:  repoURL,
		RepoPath: repoPath,

		main: main{
			fs:   fs,
			path: filepath.Join(path, "main"),
		},
		modules: modules{
			fs:   fs,
			path: filepath.Join(path, "modules"),
		},
		vtd: VTD{
			fs:   fs,
			path: filepath.Join(path, "vtd"),
		},
	}
	return ctd, nil
}

// ListMainGlobalFiles returns all TF files from the global main folder
func (ctd *CTD) ListMainGlobalFiles() ([]string, error) {
	return afero.Glob(ctd.main.fs, filepath.Join(ctd.main.globalPath(), "/*.tf"))
}

// ListMainUserFiles returns all TF files from the user main folder
func (ctd *CTD) ListMainUserFiles(user string) ([]string, error) {
	return afero.Glob(ctd.main.fs, filepath.Join(ctd.main.userPath(user), "/*.tf"))
}

func (ctd *CTD) ListModules() ([]string, error) {
	// TODO: modules could be obtained from other repos
	return afero.Glob(ctd.modules.fs, filepath.Join(ctd.modules.path, "/*"))
}

func (m *main) globalPath() string {
	return filepath.Join(m.path, "global")
}

func (m *main) userPath(user string) string {
	return filepath.Join(m.path, "user", user)
}
