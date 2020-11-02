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
	CreateUsercomponent(user string) error
	DeleteUsercomponent(user string) error
	ListUsercomponents() ([]string, error)

	CreatePluginGlobal(name string, repo string, repoPath string) error
	DeletePluginGlobal(name string) error
	ListPluginsGlobal() ([]string, error)

	CreatePluginUser(name string, user string) error
	DeletePluginUser(name string, user string) error
	ListPluginsUser(user string) ([]string, error)

	GetFlavourGlobal() (string, error)
	SetFlavourGlobal(flavour string) error

	GetFlavourUser(user string) (string, error)
	SetFlavourUser(flavour string, user string) error

	GenerateWorkdirGlobal() (string, error)
	GenerateWorkdirUser(user string) (string, error)

	GenerateVariablesGlobal() ([]string, error)
	GenerateVariablesUser(user string) ([]string, error)

	GetVariableFilepath(kind string, plugin string, user string) (string, error)
	ReadVariableFile(kind string, plugin string, user string) (string, error)

	Push(message string) error
	Pull() error

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
	Plugins [](*CTD)

	Workdir *Workdir
}

// CreateUsercomponent creates a new user component for the deployment,
// calling the respective methods con Vars and State objects
func (d *DeploymentImpl) CreateUsercomponent(user string) error {
	err := d.Vars.CreateUsercomponent(user)
	if err != nil {
		return err
	}

	err = d.State.CreateUsercomponent(user)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUsercomponent deletes an user component for the deployment,
// calling the respective methods con Vars and State objects
func (d *DeploymentImpl) DeleteUsercomponent(user string) error {
	err := d.Vars.DeleteUsercomponent(user)
	if err != nil {
		return err
	}

	err = d.State.DeleteUsercomponent(user)
	if err != nil {
		return err
	}

	return nil
}

// ListUsercomponents returns the user component names list
func (d *DeploymentImpl) ListUsercomponents() ([]string, error) {
	return d.Vars.Metadata.ListUsercomponents()
}

// CreatePluginGlobal adds a plugin to the global component, cloning its repo.
func (d *DeploymentImpl) CreatePluginGlobal(name string, repo string, repoPath string) error {
	// TODO: version and commit
	pluginPath := d.getPluginPath(name)

	err := d.fs.MkdirAll(pluginPath, 0755)
	if err != nil {
		return errors.Wrap(err, "couldn't create directory")
	}

	err = d.Vars.Metadata.CreateGlobalPlugin(name, repo, repoPath, "master", "")
	if err != nil {
		return err
	}

	plugin := NewCTD(d.fs, pluginPath, name, repo, repoPath)

	err = plugin.Clone()
	if err != nil {
		// Rollback metadata registration
		d.Vars.Metadata.DeleteGlobalPlugin(name)
		return err
	}

	d.Plugins = append(d.Plugins, plugin)
	return nil
}

// DeletePluginGlobal removes the plugin from the global component.
// TODO: any usercomponent can't have the plugin installed, must be checked
func (d *DeploymentImpl) DeletePluginGlobal(name string) error {
	err := d.fs.RemoveAll(d.getPluginPath(name))
	if err != nil {
		return errors.Wrap(err, "couldn't remove dir recursively")
	}

	return d.Vars.Metadata.DeleteGlobalPlugin(name)
}

// ListPluginsGlobal returns a list with the names of the plugins added
// to the global component
func (d *DeploymentImpl) ListPluginsGlobal() ([]string, error) {
	return d.Vars.Metadata.ListGlobalPlugins()
}

// CreatePluginUser adds a plugin to the user component. Plugin must have been
// added to de global component
func (d *DeploymentImpl) CreatePluginUser(name string, user string) error {
	return d.Vars.Metadata.CreateUserPlugin(name, user)
}

// DeletePluginUser removes the plugin from the specified user component
func (d *DeploymentImpl) DeletePluginUser(name string, user string) error {
	return d.Vars.Metadata.DeleteUserPlugin(name, user)
}

// ListPluginsUser returns a list with the names of the plugins added to
// the specified user component
func (d *DeploymentImpl) ListPluginsUser(user string) ([]string, error) {
	return d.Vars.Metadata.ListUserPlugins(user)
}

// GetFlavourGlobal returns the name of the configured flavour for the global component
func (d *DeploymentImpl) GetFlavourGlobal() (string, error) {
	return d.Vars.Metadata.GetGlobalFlavour()
}

// SetFlavourGlobal configures the flavour for the global component
func (d *DeploymentImpl) SetFlavourGlobal(flavour string) error {
	return d.Vars.Metadata.SetGlobalFlavour(flavour)
}

// GetFlavourUser returns the name of the configured flavour for the specified user component
func (d *DeploymentImpl) GetFlavourUser(user string) (string, error) {
	return d.Vars.Metadata.GetUserFlavour(user)
}

// SetFlavourUser configures the flavour for the specified user component
func (d *DeploymentImpl) SetFlavourUser(flavour string, user string) error {
	return d.Vars.Metadata.SetUserFlavour(flavour, user)
}

// GenerateWorkdirGlobal combines deployment CTDs (main and plugins) to generate
// the CTD to be applied by terraform. Returns main path where terraform must
// be executed.
func (d *DeploymentImpl) GenerateWorkdirGlobal() (string, error) {
	err := d.Workdir.GenerateGlobal()
	if err != nil {
		return "", err
	}
	return d.Workdir.mainGlobalPath(), nil
}

// GenerateWorkdirUser combines deployment CTDs (main and plugins) to generate
// the CTD to be applied by terraform
func (d *DeploymentImpl) GenerateWorkdirUser(user string) (string, error) {
	ok, err := d.Vars.Metadata.CheckUsercomponent(user)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.Errorf("user component %s doesn't exist", user)
	}

	err = d.Workdir.GenerateUser(user)
	if err != nil {
		return "", err
	}
	return d.Workdir.mainUserPath(user), nil
}

// GenerateVariablesGlobal reads the variable files from the VTDs for the global component,
// and copy them to the storage repository. Also returns a list with the filepaths and the order
// that must be used when passing the variables to Terraform.
func (d *DeploymentImpl) GenerateVariablesGlobal() ([]string, error) {
	return d.Vars.GenerateGlobal()
}

// GenerateVariablesUser reads the variable files from the VTDs for the specified user component,
// and copy them to the storage repository. Also returns a list with the filepaths and the order
// that must be used when passing the variables to Terraform.
func (d *DeploymentImpl) GenerateVariablesUser(user string) ([]string, error) {
	return d.Vars.GenerateUser(user)
}

// GetVariableFilepath returns the path where is the variable file copied to the storage repository
// for an specified kind (config, flavour or static), plugin and user component.
// Use empty string ("") on plugin parameter to obtain the base variable file and also on the
// user parameter, to obtain the global component variable file.
func (d *DeploymentImpl) GetVariableFilepath(kind string, plugin string, user string) (string, error) {
	return d.Vars.GetVariableFilepath(kind, plugin, user)
}

// ReadVariableFile returns the content of the variable file copied to the storage repository
// for an specified kind (config, flavour or static), plugin and user component.
// Use empty string ("") on plugin parameter to obtain the base variable file and also on the
// user parameter, to obtain the global component variable file.
func (d *DeploymentImpl) ReadVariableFile(kind string, plugin string, user string) (string, error) {
	filepath, err := d.GetVariableFilepath(kind, plugin, user)
	if err != nil {
		return "", err
	}

	bytes, err := afero.ReadFile(d.fs, filepath)
	if err != nil {
		return "", errors.Wrap(err, "couldn't read file")
	}

	return string(bytes), nil
}

// Push uploads vars and state to the respective repositories
func (d *DeploymentImpl) Push(message string) error {
	err := d.State.Push(message)
	if err != nil {
		return err
	}

	err = d.Vars.Push(message)
	if err != nil {
		return err
	}

	return nil
}

// Pull downloads vars and state from the respective repositories
func (d *DeploymentImpl) Pull() error {
	err := d.State.Pull()
	if err != nil {
		return err
	}

	err = d.Vars.Pull()
	if err != nil {
		return err
	}

	return nil
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
		return errors.Wrap(err, "couldn't remove dir recursively")
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

	d.Base = NewCTD(d.fs, basePath, "", d.CodeRepoURL(), d.CodeRepoPath())

	for _, plugin := range d.Vars.Metadata.Plugins {
		pluginPath := d.getPluginPath(plugin.Name)

		err := d.fs.MkdirAll(pluginPath, 0755)
		if err != nil {
			return errors.Wrap(err, "couldn't create directory")
		}

		d.Plugins = append(d.Plugins, NewCTD(d.fs, pluginPath, plugin.Name, plugin.Repo, plugin.RepoPath))
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

func (d *DeploymentImpl) getPluginPath(name string) string {
	return filepath.Join(d.path, "code", "plugins", name)
}

func (d *DeploymentImpl) getPluginByName(name string) (*CTD, error) {
	for _, plugin := range d.Plugins {
		if plugin.Name == name {
			return plugin, nil
		}
	}

	return nil, errors.Errorf("couldn't find plugin with name %s", name)
}
