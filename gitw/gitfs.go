package gitw

import (
	"github.com/spf13/afero"
	desfacer "gopkg.in/jfontan/go-billy-desfacer.v0"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

type gitfs struct {
	fs       afero.Fs
	path     string
	storer   storage.Storer
	worktree billy.Filesystem
}

// newGitFs configures the storer and the billy filesystem required to work with git
// functions, based on an Afero fs and a repository path, returning them in a gitfs struct
func newGitFs(fs afero.Fs, path string) (gitfs, error) {
	gitfs := gitfs{
		fs:       fs,
		path:     path,
		storer:   nil,
		worktree: nil,
	}

	gitfs.worktree = desfacer.New(afero.NewBasePathFs(fs, path))

	gitdir, err := gitfs.worktree.Chroot(".git")
	if err != nil {
		return gitfs, err
	}

	gitfs.storer = filesystem.NewStorage(gitdir, cache.NewObjectLRUDefault())

	return gitfs, nil
}
