package deployment

import "github.com/spf13/afero"

type workdir struct {
	fs   afero.Fs
	path string
}

func newWorkdir(fs afero.Fs, path string) (workdir, error) {
	fs.MkdirAll(path, 0700)

	workdir := workdir{
		fs:   fs,
		path: path,
	}
	return workdir, nil
}

func (w *workdir) cleanup() error {
	err := w.fs.RemoveAll(w.path)
	if err != nil {
		return err
	}
	return nil
}
