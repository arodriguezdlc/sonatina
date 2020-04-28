package deployment

import (
	"path/filepath"

	"github.com/arodriguezdlc/sonatina/utils"
	"github.com/spf13/afero"
)

type Workdir struct {
	fs   afero.Fs
	path string

	deployment *DeploymentImpl

	CTD *CTD
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

func (w *Workdir) copyMainGlobal() error {
	fileList, err := w.calculateMainGlobalFileList()
	if err != nil {
		return err
	}

	err = w.fs.MkdirAll(w.CTD.main.globalPath(), 0755)
	if err != nil {
		return err
	}

	for _, file := range fileList {
		fileName := filepath.Base(file)
		err = utils.FileCopy(w.fs, file, filepath.Join(w.CTD.main.globalPath(), fileName))
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Workdir) copyMainUser(user string) error {
	fileList, err := w.calculateMainUserFileList(user)
	if err != nil {
		return err
	}

	err = w.fs.MkdirAll(w.CTD.main.userPath(user), 0755)
	if err != nil {
		return err
	}

	for _, file := range fileList {
		fileName := filepath.Base(file)
		err = utils.FileCopy(w.fs, file, filepath.Join(w.CTD.main.userPath(user), fileName))
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

	err = w.fs.MkdirAll(w.CTD.modules.path, 0755)
	if err != nil {
		return err
	}

	for _, module := range moduleList {
		moduleName := filepath.Base(module)
		err = utils.FileCopyRecursively(w.fs, module, filepath.Join(w.CTD.modules.path, moduleName))
		if err != nil {
			return err
		}
	}

	return nil
}

// copyVTD copies variable files from VTDs in order
// XX_YY_name.tfvars
func (w *Workdir) copyVTD() error {
	w.deployment.Base.vtd.ListStaticGlobal()
	return nil
}

func (w *Workdir) cleanGlobal() error {
	err := w.fs.RemoveAll(w.CTD.main.globalPath())
	if err != nil {
		return err
	}

	err = w.fs.RemoveAll(w.CTD.modules.path)
	if err != nil {
		return err
	}

	return nil
}

func (w *Workdir) cleanUser(user string) error {
	err := w.fs.RemoveAll(w.CTD.main.userPath(user))
	if err != nil {
		return err
	}

	err = w.fs.RemoveAll(w.CTD.modules.path)
	if err != nil {
		return err
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
