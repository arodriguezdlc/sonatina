package deployment

import (
	"path/filepath"

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
	fs   afero.Fs
	path string

	Name string

	State *State
	Vars  *Vars

	CodeRepoURL string //TODO: Delete and save on CTD instead
	Base        *CTD
	Plugins     [](*CTD)

	Workdir *Workdir
}

// NewDeployment creates and initializes a new Deployment object
func NewDeployment(name string, storageRepoURL string, codeRepoURL string, fs afero.Fs, deploymentPath string) (Deployment, error) {
	deploy := &DeploymentImpl{
		Name:        name,
		fs:          fs,
		path:        deploymentPath,
		Vars:        nil,
		State:       nil,
		CodeRepoURL: codeRepoURL,
		Base:        nil,
		Plugins:     nil,
	}

	//TODO: initialize code repository using codeRepoURL
	err := deploy.initialize(storageRepoURL)
	if err != nil {
		return nil, err
	}

	return deploy, nil
}

func (d *DeploymentImpl) newWorkdir() (Workdir, error) {
	path := filepath.Join(d.path, "workdir")

	workdir := Workdir{
		fs:   d.fs,
		path: path,

		deployment: d,

		CTD: nil,
	}

	err := d.fs.MkdirAll(path, 0700)
	if err != nil {
		return workdir, err
	}

	return workdir, nil
}

// Purge removes all local files related to a deployment
func (d *DeploymentImpl) Purge() error {
	logrus.Debugln("deploymentImpl.Purge: recursive deletion on deployment path " + d.path)
	return d.fs.RemoveAll(d.path)
}

func (d *DeploymentImpl) initialize(storageRepoURL string) error {
	var err error
	var vars *Vars
	var state *State

	//Create deployment directory (idempotent operation)
	if err = d.fs.MkdirAll(d.path, 0700); err != nil {
		return err
	}

	if vars, err = NewVars(d.fs, filepath.Join(d.path, "variables"), storageRepoURL); err != nil {
		d.rollbackInitialize()
		return err
	}

	if state, err = NewState(d.fs, filepath.Join(d.path, "/state"), storageRepoURL); err != nil {
		d.rollbackInitialize()
		return err
	}

	d.Vars = vars
	d.State = state

	return nil
}

func (d *DeploymentImpl) rollbackInitialize() error {
	err := d.Purge()
	if err != nil {
		logrus.Errorln("deploymentImpl.rollbackInitialize: Error executing rollback for deployment " + d.Name)
	}
	return err
}
