package deployment

import (
	"path/filepath"

	"github.com/arodriguezdlc/sonatina/gitw"
	"github.com/arodriguezdlc/sonatina/utils"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

const stateBranch string = "state"

// State manages terraform state
type State struct {
	fs   afero.Fs
	path string
	gitw *gitw.Command

	RepoURL string
}

func (s *State) FilePathGlobal() string {
	return filepath.Join(s.path, "global", "terraform.tfstate")
}

func (s *State) FilePathUser(user string) string {
	return filepath.Join(s.path, "user", user, "terraform.tfstate")
}

// CreateUsercomponent adds a new user component to state,
// initializing the directory tree
func (s *State) CreateUsercomponent(user string) error {
	err := s.fs.MkdirAll(s.UsercomponentPath(user), 0755)
	if err != nil {
		return errors.Wrap(err, "couldn't create directory")
	}

	return nil
}

// DeleteUsercomponent deletes an user component from state and
// performs a directory cleanup
func (s *State) DeleteUsercomponent(user string) error {
	err := s.fs.RemoveAll(s.UsercomponentPath(user))
	if err != nil {
		return errors.Wrap(err, "couldn't remove dir recursively")
	}

	return nil
}

// UsercomponentPath returns de state directory path for a
// specified user
func (s *State) UsercomponentPath(user string) string {
	return filepath.Join(s.path, "user", user)
}

// Pull method retrieves terraform state information from git repository
func (s *State) Pull() error {
	return s.gitw.Pull("origin", stateBranch)
}

// Push stores terraform state information on git repository
func (s *State) Push(message string) error {
	err := s.gitw.AddGlob(".")
	if err != nil {
		return err
	}

	err = s.gitw.Commit(message)
	if err != nil {
		return err
	}

	err = s.gitw.Push("origin", stateBranch)
	if err != nil {
		return err
	}

	return nil
}

func (d *DeploymentImpl) getState(repoURL string) error {
	state, err := newState(d.fs, d.path, repoURL)
	if err != nil {
		return err
	}

	d.State = state
	return nil
}

func (d *DeploymentImpl) cloneState(repoURL string) error {
	state, err := newState(d.fs, d.path, repoURL)
	if err != nil {
		return err
	}

	err = state.gitw.CloneBranch(repoURL, stateBranch)
	if err != nil {
		return err
	}

	d.State = state
	return nil
}

func (d *DeploymentImpl) createState(repoURL string) error {
	state, err := newState(d.fs, d.path, repoURL)
	if err != nil {
		return err
	}

	err = state.gitw.Init()
	if err != nil {
		return err
	}

	err = state.gitw.RemoteAdd("origin", repoURL)
	if err != nil {
		return err
	}

	for _, subdir := range []string{"global", "user"} {
		err := utils.NewDirectoryWithKeep(state.fs, filepath.Join(state.path, subdir))
		if err != nil {
			return err
		}
	}

	err = state.gitw.AddGlob(".")
	if err != nil {
		return err
	}

	err = state.gitw.Commit("Initial commit")
	if err != nil {
		return err
	}

	// XXX: checkout executed after first commit to avoid
	// reference errors. Must be reviewed
	err = state.gitw.CheckoutNewBranch(stateBranch)
	if err != nil {
		return err
	}

	d.State = state
	return nil
}

func newState(fs afero.Fs, deploymentPath string, repoURL string) (*State, error) {
	path := filepath.Join(deploymentPath, stateBranch)

	stateGit, err := gitw.NewCommand(fs, path)
	if err != nil {
		return nil, err
	}

	state := &State{
		fs:   fs,
		path: path,
		gitw: stateGit,

		RepoURL: repoURL,
	}

	return state, nil
}
