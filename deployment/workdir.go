package deployment

import (
	"path/filepath"

	"github.com/arodriguezdlc/sonatina/utils"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type Workdir struct {
	fs   afero.Fs
	path string

	deployment *DeploymentImpl
}

func (w *Workdir) GenerateGlobal() error {
	err := w.cleanGlobal()
	if err != nil {
		return err
	}

	err = w.copyMainGlobal()
	if err != nil {
		return err
	}

	err = w.copyModules()
	if err != nil {
		return err
	}

	return nil
}

func (w *Workdir) GenerateUser(user string) error {
	err := w.cleanUser(user)
	if err != nil {
		return err
	}

	err = w.copyMainUser(user)
	if err != nil {
		return err
	}

	err = w.copyModules()
	if err != nil {
		return err
	}

	return nil
}

func (d *DeploymentImpl) newWorkdir() error {
	path := filepath.Join(d.path, "workdir")

	workdir := &Workdir{
		fs:   d.fs,
		path: path,

		deployment: d,
	}

	err := d.fs.MkdirAll(path, 0755)
	if err != nil {
		return errors.Wrapf(err, "couldn't create directory %s", path)
	}

	d.Workdir = workdir
	return nil
}

func (w *Workdir) mainGlobalPath() string {
	return filepath.Join(w.path, "main", "global")
}

func (w *Workdir) mainUserPath(user string) string {
	return filepath.Join(w.path, "main", "user", user)
}

func (w *Workdir) modulesPath() string {
	return filepath.Join(w.path, "modules")
}

func (w *Workdir) copyMainGlobal() error {
	fileList, err := w.calculateMainGlobalFileList()
	if err != nil {
		return err
	}

	mainPath := w.mainGlobalPath()
	err = w.fs.MkdirAll(mainPath, 0755)
	if err != nil {
		return errors.Wrapf(err, "couldn't create directory %s", mainPath)
	}

	for _, src := range fileList {
		dst := filepath.Join(w.mainGlobalPath(), filepath.Base(src))
		err = utils.FileCopy(w.fs, src, dst)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Workdir) copyMainUser(user string) error {
	fileList, err := w.calculateMainUserFileList()
	if err != nil {
		return err
	}

	mainPath := w.mainUserPath(user)
	err = w.fs.MkdirAll(mainPath, 0755)
	if err != nil {
		return errors.Wrapf(err, "couldn't create directory %s", mainPath)
	}

	for _, src := range fileList {
		dst := filepath.Join(w.mainUserPath(user), filepath.Base(src))
		err = utils.FileCopy(w.fs, src, dst)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Workdir) copyModules() error {
	moduleList, err := w.calculateModuleList()
	if err != nil {
		return err
	}

	err = w.fs.MkdirAll(w.modulesPath(), 0755)
	if err != nil {
		return errors.Wrap(err, "couldn't create directory")
	}

	for _, src := range moduleList {
		dst := filepath.Join(w.modulesPath(), filepath.Base(src))

		err = w.fs.MkdirAll(dst, 0755)
		if err != nil {
			return errors.Wrap(err, "couldn't create directory")
		}

		err = utils.FileCopyRecursively(w.fs, src, dst)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Workdir) cleanGlobal() error {
	path := w.mainGlobalPath()
	err := w.fs.RemoveAll(path)
	if err != nil {
		return errors.Wrap(err, "couldn't remove dir recursively")
	}

	path = w.modulesPath()
	err = w.fs.RemoveAll(path)
	if err != nil {
		return errors.Wrap(err, "couldn't remove dir recursively")
	}

	return nil
}

func (w *Workdir) cleanUser(user string) error {
	path := w.mainUserPath(user)
	err := w.fs.RemoveAll(path)
	if err != nil {
		return errors.Wrap(err, "couldn't remove dir recursively")
	}

	path = w.modulesPath()
	err = w.fs.RemoveAll(path)
	if err != nil {
		return errors.Wrap(err, "couldn't remove dir recursively")
	}

	return nil
}

func (w *Workdir) calculateMainGlobalFileList() ([]string, error) {
	files, err := w.deployment.Base.ListMainGlobalFiles()
	if err != nil {
		return nil, err
	}

	for _, plugin := range w.deployment.Plugins {
		pluginFiles, err := plugin.ListMainGlobalFiles()
		if err != nil {
			return nil, err
		}
		files = append(files, pluginFiles...)
	}

	return files, nil
}

func (w *Workdir) calculateMainUserFileList() ([]string, error) {
	files, err := w.deployment.Base.ListMainUserFiles()
	if err != nil {
		return nil, err
	}

	for _, plugin := range w.deployment.Plugins {
		pluginFiles, err := plugin.ListMainUserFiles()
		if err != nil {
			return nil, err
		}
		files = append(files, pluginFiles...)
	}

	return files, nil
}

func (w *Workdir) calculateModuleList() ([]string, error) {
	modules, err := w.deployment.Base.ListModules()
	if err != nil {
		return nil, err
	}

	for _, plugin := range w.deployment.Plugins {
		pluginModules, err := plugin.ListModules()
		if err != nil {
			return nil, err
		}
		modules = append(modules, pluginModules...)
	}

	return modules, nil
}
