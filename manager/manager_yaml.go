package manager

import (
	"errors"

	"github.com/arodriguezdlc/sonatina/utils"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"

	"github.com/arodriguezdlc/sonatina/deployment"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// ManagerYaml manages deployments saving metadata on YAML file
type ManagerYaml struct {
	Fs       afero.Fs
	FilePath string
}

type DeploymentItem struct {
	StorageRepoURI string `yaml:"storage_repo_uri"`
	CodeRepoURI    string `yaml:"code_repo_uri"`
}

type DeploymentMap map[string]DeploymentItem

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

	if err = fs.MkdirAll(path, 0755); err != nil {
		return manager, err
	}

	filepath = path + "/" + deploymentsFilename

	// Create file if doesn't exist
	if err = utils.NewFileIfNotExist(filepath, fs); err != nil {
		return manager, err
	}

	manager = ManagerYaml{
		Fs:       fs,
		FilePath: filepath,
	}

	return manager, err
}

// List returns a string slice with deployment names
func (m ManagerYaml) List() ([]string, error) {
	var d DeploymentMap
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
	var dm DeploymentMap
	var di DeploymentItem
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
		m.Fs,
		viper.GetString("DeploymentsPath"))
	if err != nil {
		return nil, err
	}

	return deploy, nil
}

// Create instantiates a new Deployment object and returns it
func (m ManagerYaml) Create(name string, storageRepoURI string, codeRepoURI string) (deployment.Deployment, error) {

	return deployment.DeploymentImpl{}, nil
}

// Delete removes the deployment from the list
func (m ManagerYaml) Delete(name string) error {
	log.Infoln("Delete " + name + " deployment")
	return nil
}

func (m ManagerYaml) read() (DeploymentMap, error) {
	var d DeploymentMap
	data, err := afero.ReadFile(m.Fs, m.FilePath)
	if err != nil {
		log.Errorln(err)
	} else {
		err = yaml.Unmarshal(data, &d)
		if err != nil {
			log.Errorln(err)
		} else {
			log.Debugln(d)
		}
	}
	return d, err
}

func (m ManagerYaml) save(d DeploymentMap) {
	data, err := yaml.Marshal(d)
	if err != nil {
		log.Fatalln(err)
	}
	err = afero.WriteFile(m.Fs, m.FilePath, data, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

// TO DO: Check that readed DeploymentMap is correct
