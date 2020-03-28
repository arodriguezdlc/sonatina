package manager

import (
	"encoding/json"

	"github.com/arodriguezdlc/sonatina/utils"

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

	err = fs.MkdirAll(path, 0700)
	if err != nil {
		return nil, err
	}

	filepath := path + "/" + deploymentsFilename

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
// added previously
func (m *managerJSON) Get(name string) (deployment.Deployment, error) {
	var dm deploymentMap
	var di deploymentItem
	var deploy deployment.Deployment
	var err error
	var ok bool

	if dm, err = m.read(); err != nil {
		return nil, err
	}

	if di, ok = dm[name]; !ok {
		return nil, DeploymentDoNotExistsError{name}
	}

	deploy, err = deployment.NewDeployment(
		name,
		di.StorageRepoURI,
		di.CodeRepoURI,
		m.fs,
		m.deploymentsPath+"/"+name)
	if err != nil {
		return nil, err
	}

	return deploy, nil
}

// Add instantiates a new Deployment object and returns it
func (m *managerJSON) Add(name string, storageRepoURI string, codeRepoURI string) (deployment.Deployment, error) {
	var dm deploymentMap
	var di deploymentItem
	var deploy deployment.Deployment
	var err error
	var ok bool

	if dm, err = m.read(); err != nil {
		return nil, err
	}

	if _, ok = dm[name]; ok {
		return nil, DeploymentAlreadyExistsError{name}
	}

	if deploy, err = deployment.NewDeployment(name, storageRepoURI, codeRepoURI, m.fs, m.deploymentsPath+"/"+name); err != nil {
		return nil, err
	}

	di = deploymentItem{
		StorageRepoURI: storageRepoURI,
		CodeRepoURI:    codeRepoURI,
	}
	m.add(name, di, &dm)
	if err = m.save(dm); err != nil {
		return nil, err
	}

	return deploy, nil
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

func (m *managerJSON) exist(name string) (bool, error) {
	deploys, err := m.read()
	if err != nil {
		return false, err
	}

	_, ok := deploys[name]
	return ok, err
}

func (m *managerJSON) get(name string) (deploymentMap, error) {
	var (
		result  deploymentMap
		deploys deploymentMap
		item    deploymentItem
		err     error
		ok      bool
	)

	if deploys, err = m.read(); err != nil {
		return nil, err
	}

	if item, ok = deploys[name]; !ok {
		return nil, DeploymentDoNotExistsError{name}
	}

	result[name] = item
	return result, err
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
		return d, err
	}

	err = json.Unmarshal(data, &d)
	if err != nil {
		logrus.Errorln(err)
	}

	return d, err
}

func (m *managerJSON) save(d deploymentMap) error {
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}

	err = afero.WriteFile(m.fs, m.filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// TODO: Check that readed DeploymentMap is correct

// TODO: some functions have to read json more than once. It could be optimized
