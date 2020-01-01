package utils

import (
	"github.com/spf13/afero"
	desfacer "gopkg.in/jfontan/go-billy-desfacer.v0"
	"gopkg.in/src-d/go-billy.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

// GitFs returns the storer and the billy filesystem required to work with git
// functions, based on an Afero fs and a repository path.
func GitFs(fs afero.Fs, path string) (storage.Storer, billy.Filesystem, error) {

	workdir := desfacer.New(afero.NewBasePathFs(fs, path))

	gitdir, err := workdir.Chroot(".git")
	if err != nil {
		return nil, nil, err
	}

	storer := filesystem.NewStorage(gitdir, cache.NewObjectLRUDefault())

	return storer, workdir, nil
}
