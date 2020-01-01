package deployment

import (
	"github.com/arodriguezdlc/sonatina/utils"
	"github.com/spf13/afero"
	"gopkg.in/src-d/go-git.v4"
)

// VarsGit implements Vars interface
type VarsGit struct {
	fs      afero.Fs
	path    string
	repoURL string
}

// NewVarsGit creates and initializes a new VarsGit object
func NewVarsGit(fs afero.Fs, path string, repoURL string) (Vars, error) {
	var err error

	storer, workdir, err := utils.GitFs(fs, path)
	if err != nil {
		return nil, err
	}

	_, err = git.Clone(storer, workdir, &git.CloneOptions{
		URL:           repoURL,
		ReferenceName: "refs/heads/variables",
		SingleBranch:  true,
	})
	if err != nil {
		return nil, err
	}

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
