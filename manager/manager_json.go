package manager

// TODO: Check that readed DeploymentMap is correct
// TODO: some functions have to read json more than once. It could be optimized

import (
	"encoding/json"
	"path/filepath"

	"github.com/arodriguezdlc/sonatina/utils"
	"github.com/pkg/errors"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"

	"github.com/arodriguezdlc/sonatina/deployment"
	"github.com/spf13/afero"
)

// managerJSON manages deployments saving metadata on JSON file
type managerJSON struct {
	fs              afero.Fs
	filePath        string
	deploymentsPath string
}

type deploymentItem struct {
	StorageRepoURI string `json:"storage_repo_uri"`
	CodeRepoURI    string `json:"code_repo_uri"`
}
type deploymentMap map[string]deploymentItem

// newManagerJSON creates and initializes a new ManagerJSON object
func newManagerJSON(fs afero.Fs, deploymentsPath string, deploymentsFilename string) (Manager, error) {
	//Initialize deployments directory
	path, err := homedir.Expand(deploymentsPath)
	if err != nil {
		return nil, err
	}

	err = fs.MkdirAll(path, 0755)
	if err != nil {
		return nil, err
	}

	filepath := filepath.Join(path, deploymentsFilename)

	manager = &managerJSON{
		fs:              fs,
		filePath:        filepath,
		deploymentsPath: path,
	}

	// Create file if doesn't exist
	if err = utils.NewFileWithContentIfNotExist(fs, filepath, "{}"); err != nil {
		return nil, err
	}

	return manager, err
}

// List returns a string slice with deployment names
func (m *managerJSON) List() ([]string, error) {
	var d deploymentMap
	var err error

	if d, err = m.read(); err != nil {
		return nil, err
	}

	keys := []string{}
	for k := range d {
		keys = append(keys, k)
	}

	return keys, nil
}

// Get instantiates and returns a Deployment object with the specified name, that have been
// initialized or cloned previously
func (m *managerJSON) Get(name string) (deployment.Deployment, error) {
	dm, err := m.read()
	if err != nil {
		return nil, err
	}

	di, ok := dm[name]
	if !ok {
		return nil, DeploymentDoNotExistsError{name}
	}

	deploy, err := deployment.Get(
		name,
		di.StorageRepoURI,
		m.fs,
		filepath.Join(m.deploymentsPath, name),
	)
	if err != nil {
		return nil, err
	}

	return deploy, nil
}

// Clone downloads deployment information from the storage repo initializes all the
// background file structure of a sonatina deployment
func (m *managerJSON) Clone(name string, storageRepoURI string, codeRepoURI string) error {
	dm, err := m.read()
	if err != nil {
		return err
	}

	if _, ok := dm[name]; ok {
		return DeploymentAlreadyExistsError{name}
	}

	err = deployment.Clone(name, storageRepoURI, m.fs, filepath.Join(m.deploymentsPath, name))
	if err != nil {
		return err
	}

	di := deploymentItem{
		StorageRepoURI: storageRepoURI,
		CodeRepoURI:    codeRepoURI,
	}
	m.add(name, di, &dm)
	if err = m.save(dm); err != nil {
		return err
	}

	return nil
}

func (m *managerJSON) Create(name string, storageRepoURI string, codeRepoURI string, codeRepoPath string,
	terraformVersion string, flavour string) error {

	dm, err := m.read()
	if err != nil {
		return err
	}

	if _, ok := dm[name]; ok {
		return DeploymentAlreadyExistsError{name}
	}

	err = deployment.Create(name, storageRepoURI, codeRepoURI, codeRepoPath,
		terraformVersion, flavour, m.fs, filepath.Join(m.deploymentsPath, name))
	if err != nil {
		return err
	}

	di := deploymentItem{
		StorageRepoURI: storageRepoURI,
		CodeRepoURI:    codeRepoURI,
	}
	m.add(name, di, &dm)
	if err = m.save(dm); err != nil {
		return err
	}

	return nil
}

// Delete removes the deployment from the list
func (m *managerJSON) Delete(name string) error {
	logrus.Infoln("Delete " + name + " deployment")

	deploy, err := m.Get(name)
	if err != nil {
		return err
	}

	if err = deploy.Purge(); err != nil {
		return err
	}

	err = m.delete(name)
	return err
}

func (m *managerJSON) delete(name string) error {
	deploys, err := m.read()
	if err != nil {
		return err
	}

	if _, ok := deploys[name]; !ok {
		return DeploymentDoNotExistsError{name}
	}
	delete(deploys, name)
	return m.save(deploys)
}

func (m *managerJSON) add(name string, di deploymentItem, dm *deploymentMap) {
	if len(*dm) < 1 {
		*dm = deploymentMap{name: di}
	} else {
		(*dm)[name] = di
	}
}

func (m *managerJSON) read() (deploymentMap, error) {
	var d deploymentMap

	data, err := afero.ReadFile(m.fs, m.filePath)
	if err != nil {
		return d, errors.Wrap(err, "couldn't read deployments file")
	}

	err = json.Unmarshal(data, &d)
	if err != nil {
		return d, errors.Wrap(err, "couldn't unmarshal deployments json file")
	}

	return d, err
}

func (m *managerJSON) save(d deploymentMap) error {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return errors.Wrap(err, "couldn't marshal deployemnts json file")
	}

	err = afero.WriteFile(m.fs, m.filePath, data, 0644)
	if err != nil {
		return errors.Wrap(err, "couldn't write deployments file")
	}

	return nil
}
