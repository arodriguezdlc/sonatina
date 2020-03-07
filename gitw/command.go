package gitw

import (
	"time"

	"github.com/spf13/afero"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
)

type Command struct {
	gitfs gitfs
}

// NewCommand is a constructor for a gitw Command object, used to execute
// equivalents to git commands over a specific repository.
// Returns a Command struct
func NewCommand(fs afero.Fs, path string) (Command, error) {
	var err error

	command := Command{
		gitfs: gitfs{},
	}

	command.gitfs, err = newGitFs(fs, path)
	if err != nil {
		return command, err
	}

	return command, nil
}

// Clone executes a `git clone` equivalent.
func (c Command) Clone(repoURL string) error {
	_, err := git.Clone(c.gitfs.storer, c.gitfs.worktree, &git.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		c.cloneRollback()
	}
	return err
}

// CloneBranch executes a `git clone`, but obtaining only specified branch.
func (c Command) CloneBranch(repoURL string, branch string) error {
	_, err := git.Clone(c.gitfs.storer, c.gitfs.worktree, &git.CloneOptions{
		URL:           repoURL,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
	})
	if err != nil {
		c.cloneRollback()
	}
	return err
}

// CloneOrInit executes a `git clone`, or a `git init` if repository doesn't exist
// or have not been initialized.
// XXX: Maybe should be moved to other package
func (c Command) CloneOrInit(repoURL string, branch string) (string, error) {
	operation := "clone"
	err := c.CloneBranch(repoURL, branch)

	if err == git.ErrInvalidReference ||
		err == transport.ErrEmptyRemoteRepository {

		operation = "init"
		_, err = c.Init()
	}

	return operation, err
}

func (c Command) Init() (*git.Repository, error) {
	return git.Init(c.gitfs.storer, c.gitfs.worktree)
}

func (c Command) CheckoutNewBranch(branch string) error {
	repo, err := c.open()
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	head, err := repo.Head()
	if err != nil {
		return err
	}

	return worktree.Checkout(&git.CheckoutOptions{
		Create: true,
		Branch: plumbing.NewBranchReferenceName(branch),
		Hash:   head.Hash(),
	})
}

func (c Command) AddGlob(pattern string) error {
	worktree, err := c.worktree()
	if err != nil {
		return err
	}

	return worktree.AddGlob(pattern)
}

func (c Command) Commit(msg string) error {
	worktree, err := c.worktree()
	if err != nil {
		return err
	}

	_, err = worktree.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe", //TODO: integrate author
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})

	return err
}

// Private

func (c Command) open() (*git.Repository, error) {
	return git.Open(c.gitfs.storer, c.gitfs.worktree)
}

func (c Command) worktree() (*git.Worktree, error) {
	repo, err := c.open()
	if err != nil {
		return nil, err
	}

	return repo.Worktree()
}

func (c Command) cloneRollback() error {
	return c.gitfs.fs.RemoveAll(c.gitfs.path)
}
