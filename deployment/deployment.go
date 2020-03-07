package deployment

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

//Deployment object contains all information about a deployment
// TODO: improve explanation and interface definition
type Deployment interface {
	Purge() error
}

// DeploymentImpl implements Deployment interface
type DeploymentImpl struct {
	Name           string
	Fs             afero.Fs
	Path           string
	CodeRepoURL    string
	StorageRepoURL string
	State          State
	Vars           Vars
	/*variables (vtd)
	state
	htd
	workingDirectory*/
}

// Purge removes all local files related to a deployment
func (d DeploymentImpl) Purge() error {
	logrus.Debugln("deploymentImpl.Purge: recursive deletion on deployment path " + d.Path)
	return d.Fs.RemoveAll(d.Path)
}

func (d DeploymentImpl) initialize() error {
	var err error
	var vars Vars
	var state State

	//Create deployment directory (idempotent operation)
	if err = d.Fs.MkdirAll(d.Path, 0700); err != nil {
		return err
	}

	if vars, err = NewVarsGit(d.Fs, d.Path+"/variables", d.StorageRepoURL); err != nil {
		d.rollbackInitialize()
		return err
	}

	if state, err = NewStateGit(d.Fs, d.Path+"/state", d.StorageRepoURL); err != nil {
		d.rollbackInitialize()
		return err
	}

	d.Vars = vars
	d.State = state

	return nil
}

func (d DeploymentImpl) rollbackInitialize() error {
	err := d.Purge()
	if err != nil {
		logrus.Errorln("deploymentImpl.rollbackInitialize: Error executing rollback for deployment " + d.Name)
	}
	return err
}

// NewDeployment creates and initializes a new Deployment object
func NewDeployment(name string, storageRepoURL string, codeRepoURL string, fs afero.Fs, deploymentPath string) (Deployment, error) {
	deploy := DeploymentImpl{
		Name:           name,
		Fs:             fs,
		Path:           deploymentPath,
		CodeRepoURL:    codeRepoURL,
		StorageRepoURL: storageRepoURL,
		Vars:           nil,
		State:          nil,
	}

	err := deploy.initialize()
	if err != nil {
		return nil, err
	}

	return deploy, nil
}
