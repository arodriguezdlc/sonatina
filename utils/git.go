package utils

import (
	"github.com/spf13/afero"
	desfacer "gopkg.in/jfontan/go-billy-desfacer.v0"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

// GitCloneOrInit clones a repository or initializes it if is empty or remote
// branch doesn't exists
func GitCloneOrInit(fs afero.Fs, path string, repoURL string, branch string) error {
	err := gitClone(fs, path, repoURL, branch)

	if err == git.ErrInvalidReference {
		return gitInit(fs, path, repoURL, branch)
	}

	return nil
}

// GitOpen is a wrapper of git.Open method to use an afero.Fs and a path
// insteado of using storer and worktree parameters, making it easier to use
func GitOpen(fs afero.Fs, path string) (*git.Repository, error) {
	storer, worktree, err := gitFs(fs, path)
	if err != nil {
		return nil, err
	}

	return git.Open(storer, worktree)
}

func gitClone(fs afero.Fs, path string, repoURL string, branch string) error {
	storer, worktree, err := gitFs(fs, path)
	if err != nil {
		return err
	}

	_, err = git.Clone(storer, worktree, &git.CloneOptions{
		URL:           repoURL,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
	})
	return err
}

func gitInit(fs afero.Fs, path string, repoURL string, branch string) error {
	storer, worktree, err := gitFs(fs, path)
	if err != nil {
		return err
	}

	repo, err := git.Init(storer, worktree)
	if err != nil {
		return err
	}

	err = repo.CreateBranch(&config.Branch{
		Name: branch,
	})
	if err != nil {
		return err
	}

	repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{repoURL},
	})
	if err != nil {
		return err
	}

	return nil
}

// gitFs returns the storer and the billy filesystem required to work with git
// functions, based on an Afero fs and a repository path.
func gitFs(fs afero.Fs, path string) (storage.Storer, billy.Filesystem, error) {

	worktree := desfacer.New(afero.NewBasePathFs(fs, path))

	gitdir, err := worktree.Chroot(".git")
	if err != nil {
		return nil, nil, err
	}

	storer := filesystem.NewStorage(gitdir, cache.NewObjectLRUDefault())

	return storer, worktree, nil
}
