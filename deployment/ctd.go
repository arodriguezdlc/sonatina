package deployment

import "github.com/spf13/afero"

// CTD represents a Code Tree Definition
type CTD struct {
	fs   afero.Fs
	path string

	RepoURL  string
	RepoPath string

	main    main
	modules modules
	vtd     vtd
}

type main struct {
	fs   afero.Fs
	path string
}

type modules struct {
	fs   afero.Fs
	path string
}

type vtd struct {
	fs   afero.Fs
	path string
}

// TODO
func NewCTD(fs afero.Fs, path string, repoURL string, repoPath string) (CTD, error) {
	ctd := CTD{
		fs:   fs,
		path: path,

		RepoURL:  repoURL,
		RepoPath: repoPath,

		main: main{
			fs:   fs,
			path: path + "main",
		},
		modules: modules{
			fs:   fs,
			path: path + "modules",
		},
		vtd: vtd{
			fs:   fs,
			path: path + "vtd",
		},
	}
	return ctd, nil
}

// ListMainFilesGlobal returns all TF files from the global main folder
func (ctd *CTD) ListMainFilesGlobal() ([]string, error) {
	return afero.Glob(ctd.main.fs, ctd.main.path+"/global/*.tf")
}

// ListMainFilesUser returns all TF files from the user main folder
func (ctd *CTD) ListMainFilesUser(user string) ([]string, error) {
	return afero.Glob(ctd.main.fs, ctd.main.path+"/user/"+user+"/*.tf")
}
