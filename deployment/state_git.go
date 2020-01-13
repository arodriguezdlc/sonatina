package deployment

import (
	"github.com/arodriguezdlc/sonatina/utils"
	"github.com/spf13/afero"
)

// StateGit implements State interface
type StateGit struct {
	fs      afero.Fs
	path    string
	repoURL string
}

// NewStateGit creates and initializes a new ManagerYaml object
func NewStateGit(fs afero.Fs, path string, repoURL string) (State, error) {
	err := utils.GitCloneOrInit(fs, path, repoURL, "state")
	if err != nil {
		return nil, err
	}

	return StateGit{
		fs:      fs,
		path:    path,
		repoURL: repoURL,
	}, nil
}

// Path returns the path where state is saved
func (s StateGit) Path() string {
	return s.path
}

// Save method stores terraform state information on git repository
func (s StateGit) Save() {
	// TO DO
}

// Load method retrieves terraform state information from git repository
func (s StateGit) Load() {
	// TO DO
}

func (s StateGit) clone() {
	//_, err := git.PlainClone(s.Path)
	//git.Plai
}
