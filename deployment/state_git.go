package deployment

import (
	"github.com/arodriguezdlc/sonatina/gitw"
	"github.com/arodriguezdlc/sonatina/utils"
	"github.com/spf13/afero"
)

// StateGit implements State interface
type StateGit struct {
	fs      afero.Fs
	path    string
	repoURL string
	gitw    gitw.Command
}

// NewStateGit creates and initializes a new ManagerYaml object
func NewStateGit(fs afero.Fs, path string, repoURL string) (State, error) {
	var err error

	state := StateGit{
		fs:      fs,
		path:    path,
		repoURL: repoURL,
		gitw:    gitw.Command{},
	}

	state.gitw, err = gitw.NewCommand(fs, path)
	if err != nil {
		return nil, err
	}

	ok, err := state.isInitialized()
	if err != nil {
		return nil, err
	}
	if !ok {
		operation, err := state.gitw.CloneOrInit(repoURL, "state")
		if err != nil {
			return nil, err
		}
		if operation == "init" {
			if err = state.initialize(); err != nil {
				return nil, err
			}
		}
	}

	return state, nil
}

// Path returns the path where state is saved
func (s StateGit) Path() string {
	return s.path
}

// Save method stores terraform state information on git repository
func (s StateGit) Save() {
	// TODO
}

// Load method retrieves terraform state information from git repository
func (s StateGit) Load() {
	// TODO
}

func (s StateGit) initialize() error {
	err := s.fs.MkdirAll(s.path+"/global", 0700)
	if err != nil {
		return err
	}

	err = utils.NewFileIfNotExist(s.path+"/global/.keep", s.fs)
	if err != nil {
		return err
	}

	err = s.fs.MkdirAll(s.path+"/user", 0700)
	if err != nil {
		return err
	}

	err = utils.NewFileIfNotExist(s.path+"/user/.keep", s.fs)
	if err != nil {
		return err
	}

	err = s.gitw.AddGlob(".")
	if err != nil {
		return err
	}

	err = s.gitw.Commit("Initial commit")
	if err != nil {
		return err
	}

	return s.gitw.CheckoutNewBranch("state")
}

func (s StateGit) isInitialized() (bool, error) {
	return afero.Exists(s.fs, s.path)
}
