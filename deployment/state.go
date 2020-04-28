package deployment

import (
	"github.com/arodriguezdlc/sonatina/gitw"
	"github.com/arodriguezdlc/sonatina/utils"
	"github.com/spf13/afero"
)

// State manages terraform state
type State struct {
	fs      afero.Fs
	path    string
	repoURL string
	gitw    gitw.Command
}

// NewState creates and initializes a new State object
func NewState(fs afero.Fs, path string, repoURL string) (*State, error) {
	var err error

	state := &State{
		fs:      fs,
		path:    path,
		repoURL: repoURL,
		gitw:    gitw.Command{},
	}

	state.gitw, err = gitw.NewCommand(fs, path)
	if err != nil {
		return state, err
	}

	ok, err := state.isInitialized()
	if err != nil {
		return state, err
	}
	if !ok {
		operation, err := state.gitw.CloneOrInit(repoURL, "state")
		if err != nil {
			return state, err
		}
		if operation == "init" {
			if err = state.initialize(); err != nil {
				return state, err
			}
		}
	}

	return state, nil
}

// Save method stores terraform state information on git repository
func (s *State) Save() {
	// TODO
}

// Load method retrieves terraform state information from git repository
func (s *State) Load() {
	// TODO
}

func (s *State) initialize() error {
	err := s.fs.MkdirAll(s.path+"/global", 0755)
	if err != nil {
		return err
	}

	err = utils.NewFileIfNotExist(s.fs, s.path+"/global/.keep")
	if err != nil {
		return err
	}

	err = s.fs.MkdirAll(s.path+"/user", 0755)
	if err != nil {
		return err
	}

	err = utils.NewFileIfNotExist(s.fs, s.path+"/user/.keep")
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

func (s *State) isInitialized() (bool, error) {
	return afero.Exists(s.fs, s.path)
}
