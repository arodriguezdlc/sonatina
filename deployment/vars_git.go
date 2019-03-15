package deployment

import "github.com/spf13/afero"

// VarsGit implements Vars interface
type VarsGit struct {
	fs      afero.Fs
	path    string
	repoURL string
}

// NewVarsGit creates and initializes a new VarsGit object
func NewVarsGit(fs afero.Fs, path string, repoURL string) (Vars, error) {
	return VarsGit{
		fs:      fs,
		path:    path,
		repoURL: repoURL,
	}, nil
}

// Path returns the path where state is saved
func (v VarsGit) Path() string {
	return v.path
}

// Save method stores terraform state information on git repository
func (v VarsGit) Save() {
	// TO DO
}

// Load method retrieves terraform state information from git repository
func (v VarsGit) Load() {
	// TO DO
}
