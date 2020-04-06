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
		vtd: VTD{
			fs:   fs,
			path: path + "vtd",
		},
	}
	return ctd, nil
}

// ListMainGlobalFiles returns all TF files from the global main folder
func (ctd *CTD) ListMainGlobalFiles() ([]string, error) {
	return afero.Glob(ctd.main.fs, ctd.main.path+"/global/*.tf")
}

// ListMainUserFiles returns all TF files from the user main folder
func (ctd *CTD) ListMainUserFiles(user string) ([]string, error) {
	return afero.Glob(ctd.main.fs, ctd.main.path+"/user/"+user+"/*.tf")
}

func (ctd *CTD) ListModules() ([]string, error) {
	// TODO: modules could be obtained from other repos
	return afero.Glob(ctd.modules.fs, ctd.modules.path+"/*")
}
