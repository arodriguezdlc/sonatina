package gitw

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Command implements methods with same interface than official
// git commands but over go-git.
type Command struct {
	fs   afero.Fs
	path string
}

// NewCommand is a constructor for a gitw Command object, used to execute
// equivalents to git commands over a specific repository.
// Returns a Command struct
func NewCommand(fs afero.Fs, path string) (*Command, error) {
	// XXX: fs afero.Fs is not used, but is maintained for retrocompatibility
	// Also, we would like to use afero Fs to provide testing via ram filesystem.
	return &Command{
		fs:   fs,
		path: path,
	}, nil
}

// Clone executes a `git clone` equivalent.
func (c *Command) Clone(repoURL string) error {
	_, err := git.PlainClone(c.path, false, &git.CloneOptions{
		URL: repoURL,
	})
	return errors.Wrapf(err, "couldn't clone repository %s", repoURL)
}

// CloneBranch executes a `git clone`, but obtaining only specified branch.
func (c *Command) CloneBranch(repoURL string, branch string) error {
	_, err := git.PlainClone(c.path, false, &git.CloneOptions{
		URL:           repoURL,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
	})
	return errors.Wrapf(err, "couldn't clone repository %s", repoURL)
}

// Init executes a `git init` equivalent
func (c *Command) Init() error {
	_, err := git.PlainInit(c.path, false)
	if err != nil {
		return errors.Wrapf(err, "couldn't init new repository on %s", c.path)
	}
	return nil
}

// CheckoutNewBranch executes a `git checkout -b` equivalent
func (c *Command) CheckoutNewBranch(branch string) error {
	repo, worktree, err := c.openWithWorktree()
	if err != nil {
		return err
	}

	head, err := repo.Head()
	if err != nil {
		return errors.Wrap(err, "couldn't get head reference")
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Create: true,
		Branch: plumbing.NewBranchReferenceName(branch),
		Hash:   head.Hash(),
	})
	if err != nil {
		return errors.Wrapf(err, "couldn't checkout to new branch %s", branch)
	}

	return nil
}

// AddGlob executes a `git add` equivalent
func (c *Command) AddGlob(pattern string) error {
	worktree, err := c.worktree()
	if err != nil {
		return err
	}

	err = worktree.AddGlob(pattern)
	if err != nil {
		return errors.Wrap(err, "couldn't add to worktree")
	}

	return nil
}

// Commit executes a `git commit -m` equivalent
func (c *Command) Commit(msg string) error {
	repo, worktree, err := c.openWithWorktree()
	if err != nil {
		return err
	}

	email, err := c.getEmail(repo)
	if err != nil {
		return err
	}

	name, err := c.getUsername(repo)
	if err != nil {
		return err
	}

	_, err = worktree.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  name,
			Email: email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return errors.Wrap(err, "couldn't create commit")
	}

	return nil
}

// RemoteAdd executes a `git remote add` equivalent
func (c *Command) RemoteAdd(name string, url string) error {
	repo, err := c.open()
	if err != nil {
		return err
	}

	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: name,
		URLs: []string{url},
	})
	if err != nil {
		return errors.Wrap(err, "couldn't create remote")
	}

	return nil
}

// Pull executes a `git pull <remote> <branch>` equivalent
func (c *Command) Pull(remote string, branch string) error {
	worktree, err := c.worktree()
	if err != nil {
		return err
	}

	err = worktree.Pull(&git.PullOptions{
		RemoteName:    remote,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return errors.Wrap(err, "couldn't pull from origin")
	}

	return nil
}

// Push executes a `git push <remote> <branch>` equivalent
func (c *Command) Push(remote string, branch string) error {
	repo, err := c.open()
	if err != nil {
		return err
	}

	ref := plumbing.NewBranchReferenceName(branch)
	referenceList := append([]config.RefSpec{}, config.RefSpec(ref+":"+ref))
	err = repo.Push(&git.PushOptions{
		RemoteName: remote,
		RefSpecs:   referenceList,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	return nil
}

// TODO
func (c *Command) Reset() error {
	//worktree, err := c.worktree()
	// if err != nil {
	// 	return err
	// }

	//err = worktree.Reset(&git.ResetOptions{
	//	Commit: ,
	//	Mode: ,
	//})

	return nil
}

// Private

func (c *Command) open() (*git.Repository, error) {
	repo, err := git.PlainOpen(c.path)
	if err != nil {
		return repo, errors.Wrapf(err, "couldn't open repo on %s", c.path)
	}

	return repo, nil
}

func (c *Command) openWithWorktree() (*git.Repository, *git.Worktree, error) {
	repo, err := c.open()
	if err != nil {
		return nil, nil, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, nil, errors.Wrapf(err, "couldn't open worktree on %s", c.path)
	}

	return repo, worktree, nil
}

func (c *Command) worktree() (*git.Worktree, error) {
	_, worktree, err := c.openWithWorktree()
	return worktree, err
}

func (c *Command) getUsername(repo *git.Repository) (string, error) {
	cfg, err := repo.ConfigScoped(config.SystemScope)
	if err != nil {
		return "", errors.Wrap(err, "couldn't get username from git config")
	}

	return cfg.User.Name, nil
}

func (c *Command) getEmail(repo *git.Repository) (string, error) {
	cfg, err := repo.ConfigScoped(config.SystemScope)
	if err != nil {
		return "", errors.Wrap(err, "couldn't get email from git config")
	}

	return cfg.User.Email, nil
}
