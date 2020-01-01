package manager

import (
	"errors"

	"github.com/arodriguezdlc/sonatina/utils"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"

	"github.com/arodriguezdlc/sonatina/deployment"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

// ManagerYaml manages deployments saving metadata on YAML file
type ManagerYaml struct {
	Fs              afero.Fs
	FilePath        string
	DeploymentsPath string
}

type deploymentItem struct {
	StorageRepoURI string `yaml:"storage_repo_uri"`
	CodeRepoURI    string `yaml:"code_repo_uri"`
}
type deploymentMap map[string]deploymentItem

// NewManagerYaml creates and initializes a new ManagerYaml object
func NewManagerYaml(fs afero.Fs, deploymentsPath string, deploymentsFilename string) (ManagerYaml, error) {
	var manager ManagerYaml
	var err error
	var path string
	var filepath string

	//Initialize deployments directory
	if path, err = homedir.Expand(deploymentsPath); err != nil {
		return manager, err
	}

	if err = fs.MkdirAll(path, 0700); err != nil {
		return manager, err
	}

	filepath = path + "/" + deploymentsFilename

	// Create file if doesn't exist
	if err = utils.NewFileIfNotExist(filepath, fs); err != nil {
		return manager, err
	}

	manager = ManagerYaml{
		Fs:              fs,
		FilePath:        filepath,
		DeploymentsPath: path,
	}

	return manager, err
}

// List returns a string slice with deployment names
func (m ManagerYaml) List() ([]string, error) {
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

// Get instantiates and returns a Deployment object with the specified name
func (m ManagerYaml) Get(name string) (deployment.Deployment, error) {
	var dm deploymentMap
	var di deploymentItem
	var deploy deployment.Deployment
	var err error
	var ok bool

	if dm, err = m.read(); err != nil {
		return nil, err
	}

	if di, ok = dm[name]; !ok {
		return nil, errors.New("Can't find deployment with name [" + name + "]")
	}

	deploy, err = deployment.NewDeploymentImpl(
		name,
		di.StorageRepoURI,
		di.CodeRepoURI,
		m.Fs,
		m.DeploymentsPath+"/"+name)
	if err != nil {
		return nil, err
	}

	return deploy, nil
}

// Add instantiates a new Deployment object and returns it
func (m ManagerYaml) Add(name string, storageRepoURI string, codeRepoURI string) (deployment.Deployment, error) {
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

	if deploy, err = deployment.NewDeploymentImpl(name, storageRepoURI, codeRepoURI, m.Fs, m.DeploymentsPath+"/"+name); err != nil {
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
func (m ManagerYaml) Delete(name string) error {
	log.Infoln("Delete " + name + " deployment")
	return nil
}

func (m ManagerYaml) get(name string) (deploymentMap, error) {
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
		return nil, errors.New("Deployment doesn't exist")
	}

	result[name] = item
	return result, err
}

func (m ManagerYaml) add(name string, di deploymentItem, dm *deploymentMap) {
	if len(*dm) < 1 {
		*dm = deploymentMap{name: di}
	} else {
		(*dm)[name] = di
	}
}

func (m ManagerYaml) read() (deploymentMap, error) {
	var d deploymentMap

	data, err := afero.ReadFile(m.Fs, m.FilePath)
	if err != nil {
		return d, err
	}

	err = yaml.Unmarshal(data, &d)
	if err != nil {
		log.Errorln(err)
	}

	return d, err
}

func (m ManagerYaml) save(d deploymentMap) error {
	data, err := yaml.Marshal(d)
	if err != nil {
		return err
	}

	err = afero.WriteFile(m.Fs, m.FilePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// TODO: Check that readed DeploymentMap is correct
