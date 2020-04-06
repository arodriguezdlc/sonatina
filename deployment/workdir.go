package deployment

import "github.com/spf13/afero"

type Workdir struct {
	fs   afero.Fs
	path string

	deployment *DeploymentImpl

	CTD *CTD
	VTD *VTD
}

func newWorkdir(fs afero.Fs, path string) (Workdir, error) {
	fs.MkdirAll(path, 0700)

	workdir := Workdir{
		fs:   fs,
		path: path,
	}
	return workdir, nil
}

func (w *Workdir) GenerateGlobal() error {
	// TODO
	// Clean previous workdir
	w.cleanGlobal()
	// Calculate files to copy
	// Copy files
	return nil
}

func (w *Workdir) GenerateUser() error {
	// TODO
	return nil
}

func (w *Workdir) cleanGlobal() error {
	// TODO
	return nil
}

func (w *Workdir) cleanUser() error {
	// TODO
	return nil
}

func (w *Workdir) calculateGlobalFileList() error {
	fileList, err := w.calculateMainGlobalFileList()
	if err != nil {
		return err
	}
	// TODO: modules, vars
	return nil
}

func (w *Workdir) generateMainUser(user string) error {
	fileList, err := w.calculateMainUserFileList(user)
	if err != nil {
		return err
	}

	return nil
}

func (w *Workdir) calculateModulesFileList() error {
	// TODO
	//fileList, err := w.deployment.Base.

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

	return files, err
}

func (w *Workdir) calculateMainUserFileList(user string) ([]string, error) {
	files, err := w.deployment.Base.ListMainUserFiles(user)
	if err != nil {
		return nil, err
	}

	for _, plugin := range w.deployment.Plugins {
		pluginFiles, err := plugin.ListMainUserFiles(user)
		if err != nil {
			return nil, err
		}
		files = append(files, pluginFiles...)
	}

	return files, err
}

func (w *Workdir) cleanup() error {
	err := w.fs.RemoveAll(w.path)
	if err != nil {
		return err
	}
	return nil
}
