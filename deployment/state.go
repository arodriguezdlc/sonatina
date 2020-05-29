package deployment

import (
	"path/filepath"

	"github.com/arodriguezdlc/sonatina/gitw"
	"github.com/arodriguezdlc/sonatina/utils"
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

// Save method stores terraform state information on git repository
func (s *State) Save() {
	// TODO
}

// Load method retrieves terraform state information from git repository
func (s *State) Load() {
	// TODO
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
