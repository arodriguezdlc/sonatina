package deployment

import (
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

//Deployment object contains all information about a deployment
// TODO: improve explanation and interface definition
type Deployment interface {
	GenerateWorkdirGlobal() (string, error)
	GenerateWorkdirUser(user string) (string, error)
	GenerateVariablesGlobal() ([]string, error)
	GenerateVariablesUser(user string) ([]string, error)
	StateFilePathGlobal() string
	StateFilePathUser(user string) string
	TerraformVersion() string
	CodeRepoURL() string
	CodeRepoPath() string
	Purge() error
}

// DeploymentImpl implements Deployment interface
type DeploymentImpl struct {
	fs   afero.Fs
	path string

	Name string

	State *State
	Vars  *Vars

	Base    *CTD
	Plugins [](*CTD) // The key is the plugin name

	Workdir *Workdir
}

// GenerateWorkdirGlobal combines deployment CTDs (main and plugins) to generate
// the CTD to be applied by terraform. Returns main path where terraform must
// be executed.
func (d *DeploymentImpl) GenerateWorkdirGlobal() (string, error) {
	err := d.Workdir.GenerateGlobal()
	if err != nil {
		return "", err
	}
	return d.Workdir.CTD.main.globalPath(), nil
}

// GenerateWorkdirUser combines deployment CTDs (main and plugins) to generate
// the CTD to be applied by terraform
func (d *DeploymentImpl) GenerateWorkdirUser(user string) (string, error) {
	err := d.Workdir.GenerateUser(user)
	if err != nil {
		return "", err
	}
	return d.Workdir.CTD.main.userPath(user), nil
}

func (d *DeploymentImpl) GenerateVariablesGlobal() ([]string, error) {
	return d.Vars.GenerateGlobal()
}

func (d *DeploymentImpl) GenerateVariablesUser(user string) ([]string, error) {
	return d.Vars.GenerateUser(user)
}

func (d *DeploymentImpl) StateFilePathGlobal() string {
	return d.State.FilePathGlobal()
}

func (d *DeploymentImpl) StateFilePathUser(user string) string {
	return d.State.FilePathUser(user)
}

// TerraformVersion returns the terraform version that is being using with this
// specific deployment.
func (d *DeploymentImpl) TerraformVersion() string {
	return d.Vars.Metadata.TerraformVersion
}

// CodeRepoURL returns the URL where is the terraform code that describes
// infrastructure to be deployed in a sonatina way.
func (d *DeploymentImpl) CodeRepoURL() string {
	return d.Vars.Metadata.Repo
}

// CodeRepoPath returns the path inside the CodeRepo where is the terraform
// code that describes infrastructure to be deployed in a sonatina way.
func (d *DeploymentImpl) CodeRepoPath() string {
	return d.Vars.Metadata.RepoPath
}

// Purge removes all local files related to a deployment
func (d *DeploymentImpl) Purge() error {
	logrus.WithFields(logrus.Fields{
		"deployment": d.Name,
		"path":       d.path}).Debug("purge deployment files")

	err := d.fs.RemoveAll(d.path)
	if err != nil {
		return errors.Wrapf(err, "couldn't remove dir %s", d.path)
	}
	return nil
}

// Get creates and initializes a new Deployment object from local storage
func Get(name string, storageRepoURL string, fs afero.Fs, deploymentPath string) (Deployment, error) {
	deploy := newDeploymentImpl(name, fs, deploymentPath)

	// TODO: paralelize, using contexts to cancel operations
	err := deploy.getVars(storageRepoURL)
	if err != nil {
		return nil, err
	}

	err = deploy.getState(storageRepoURL)
	if err != nil {
		return nil, err
	}

	err = deploy.newDeploymentCTDs()
	if err != nil {
		return nil, err
	}

	err = deploy.newWorkdir()
	if err != nil {
		return nil, err
	}

	return deploy, nil
}

// Clone creates and initializes a new Deployment object that has not been created before on any repository
func Clone(name string, storageRepoURL string, fs afero.Fs, deploymentPath string) error {
	deploy := newDeploymentImpl(name, fs, deploymentPath)

	// TODO: paralelize
	err := deploy.cloneVars(storageRepoURL)
	if err != nil {
		deploy.rollbackInitialize()
		return err
	}

	err = deploy.cloneState(storageRepoURL)
	if err != nil {
		deploy.rollbackInitialize()
		return err
	}

	err = deploy.newDeploymentCTDs()
	if err != nil {
		deploy.rollbackInitialize()
		return err
	}

	err = deploy.cloneDeploymentCTDs()
	if err != nil {
		deploy.rollbackInitialize()
		return err
	}

	err = deploy.newWorkdir()
	if err != nil {
		deploy.rollbackInitialize()
		return err
	}

	return nil
}

// Create creates and initializes a new Deployment object that has not been created before on any repository
func Create(name string, storageRepoURL string, codeRepoURL string, codeRepoPath string,
	terraformVersion string, flavour string, fs afero.Fs, deploymentPath string) error {

	deploy := newDeploymentImpl(name, fs, deploymentPath)

	// TODO: paralelize
	err := deploy.createVars(storageRepoURL, terraformVersion, codeRepoURL, codeRepoPath, flavour)
	if err != nil {
		deploy.rollbackInitialize()
		return err
	}

	err = deploy.createState(storageRepoURL)
	if err != nil {
		deploy.rollbackInitialize()
		return err
	}

	err = deploy.newDeploymentCTDs()
	if err != nil {
		deploy.rollbackInitialize()
		return err
	}

	err = deploy.cloneDeploymentCTDs()
	if err != nil {
		deploy.rollbackInitialize()
		return err
	}

	err = deploy.newWorkdir()
	if err != nil {
		deploy.rollbackInitialize()
		return err
	}

	return nil
}

func newDeploymentImpl(name string, fs afero.Fs, deploymentPath string) *DeploymentImpl {
	return &DeploymentImpl{
		fs:   fs,
		path: deploymentPath,

		Name:    name,
		Vars:    nil,
		State:   nil,
		Base:    nil, //TODO
		Plugins: nil, //TODO
		Workdir: nil,
	}
}

func (d *DeploymentImpl) rollbackInitialize() error {
	return d.Purge()
}

func (d *DeploymentImpl) newDeploymentCTDs() error {
	basePath := filepath.Join(d.path, "code", "base")

	err := d.fs.MkdirAll(basePath, 0755)
	if err != nil {
		return errors.Wrap(err, "couldn't create directory")
	}

	d.Base = NewCTD(d.fs, basePath, d.CodeRepoURL(), d.CodeRepoPath())

	for _, plugin := range d.Vars.Metadata.Plugins {
		pluginPath := filepath.Join(d.path, "code", "plugins", plugin.Name)

		err := d.fs.MkdirAll(pluginPath, 0755)
		if err != nil {
			return errors.Wrap(err, "couldn't create directory")
		}

		d.Plugins = append(d.Plugins, NewCTD(d.fs, pluginPath, plugin.Repo, plugin.RepoPath))
	}

	return nil
}

func (d *DeploymentImpl) cloneDeploymentCTDs() error {
	err := d.Base.Clone()
	if err != nil {
		return err
	}

	for _, plugin := range d.Plugins {
		err = plugin.Clone()
		if err != nil {
			return err
		}
	}

	return nil
}
