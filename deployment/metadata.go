package deployment

import (
	"encoding/json"
	"path/filepath"

	"github.com/arodriguezdlc/sonatina/utils"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

const metadataFileName string = "metadata.json"

// Metadata struct is a model used for marshall/unmarshall the sonatina deployment
// metadata to/from a json file, that it's saved on the variables branch of the storage repository.
type Metadata struct {
	fs       afero.Fs
	filePath string

	TerraformVersion string                   `json:"terraform_version"`
	Repo             string                   `json:"repo"`
	RepoPath         string                   `json:"repo_path"`
	Version          string                   `json:"version"`
	Commit           string                   `json:"commit"`
	Flavour          string                   `json:"flavour"`
	UserComponents   map[string]userComponent `json:"user_components"`
	Plugins          []globalPlugin           `json:"plugins"`
}

type userComponent struct {
	Plugins []userPlugin `json:"plugins"`
	Flavour string       `json:"flavour"`
}

type globalPlugin struct {
	Name     string `json:"name"`
	Repo     string `json:"repo"`
	RepoPath string `json:"repo_path"`
	Version  string `json:"version"`
	Commit   string `json:"commit"`
}

type userPlugin struct {
	Name string `json:"name"`
}

// CreateGlobalPlugin adds to metadata a new plugin for the global component
// XXX: this method isn't thread safe
func (m *Metadata) CreateGlobalPlugin(name string, repo string, repoPath string, version string, commit string) error {
	err := m.load()
	if err != nil {
		return err
	}

	if m.globalPluginExists(name) {
		return errors.Errorf("global plugin %s already exists", name)
	}

	plugin := globalPlugin{
		Name:     name,
		Repo:     repo,
		RepoPath: repoPath,
		Version:  version,
		Commit:   commit,
	}

	m.Plugins = append(m.Plugins, plugin)

	err = m.save()
	if err != nil {
		return err
	}

	return nil
}

// DeleteGlobalPlugin deletes the specified plugin from the
// global component
// XXX: this method isn't thread safe
func (m *Metadata) DeleteGlobalPlugin(name string) error {
	err := m.load()
	if err != nil {
		return err
	}

	// TODO: check if an user component has the plugin assigned. It can't be deleted
	// until all user components have been deleted its plugin

	index, err := m.getGlobalPluginIndex(name)
	if err != nil {
		return err
	}

	m.deleteGlobalPluginWithIndex(index)

	err = m.save()
	if err != nil {
		return err
	}

	return nil
}

// ListGlobalPlugins loads metadata and list plugins added to
// the global component
func (m *Metadata) ListGlobalPlugins() ([]string, error) {
	err := m.load()
	if err != nil {
		return []string{}, err
	}

	return m.listGlobalPlugins()
}

// CreateUserPlugin adds the specified plugin to the specified user component
// XXX: this method isn't thread safe
func (m *Metadata) CreateUserPlugin(name string, user string) error {
	err := m.load()
	if err != nil {
		return err
	}

	ok, err := m.checkUsercomponent(user)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Errorf("user component %s doesn't exist", user)
	}

	if !m.globalPluginExists(name) {
		return errors.Errorf("global plugin %s doesn't exist", name)
	}

	if m.userPluginExists(name, user) {
		return errors.Errorf("user plugin %s already exists for user %s", name, user)
	}

	plugin := userPlugin{
		Name: name,
	}
	userComponent := m.UserComponents[user]
	userComponent.Plugins = append(userComponent.Plugins, plugin)
	m.UserComponents[user] = userComponent

	err = m.save()
	if err != nil {
		return err
	}

	return nil
}

// DeleteUserPlugin deletes the specified plugin from the specified
// user component
// XXX: this method isn't thread safe
func (m *Metadata) DeleteUserPlugin(name string, user string) error {
	err := m.load()
	if err != nil {
		return err
	}

	ok, err := m.checkUsercomponent(user)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Errorf("user component %s doesn't exist", user)
	}

	if !m.globalPluginExists(name) {
		return errors.Errorf("global plugin %s doesn't exist", name)
	}

	if !m.userPluginExists(name, user) {
		return errors.Errorf("user plugin %s doesn't exist for user %s", name, user)
	}

	index, err := m.getUserPluginIndex(name, user)
	if err != nil {
		return err
	}

	m.deleteUserPluginWithIndex(index, user)

	err = m.save()
	if err != nil {
		return err
	}

	return nil
}

// ListUserPlugins loads metadata and list plugins added to a
// specified user
func (m *Metadata) ListUserPlugins(user string) ([]string, error) {
	err := m.load()
	if err != nil {
		return []string{}, err
	}

	ok, err := m.checkUsercomponent(user)
	if err != nil {
		return []string{}, err
	}
	if !ok {
		return []string{}, errors.Errorf("user component %s doesn't exist", user)
	}

	return m.listUserPlugins(user)
}

// CreateUsercomponent updates metadata with a new user component
// XXX: this method isn't thread safe
func (m *Metadata) CreateUsercomponent(user string) error {
	err := m.load()
	if err != nil {
		return err
	}

	_, ok := m.UserComponents[user]
	if ok {
		return errors.Errorf("user component %s already exists", user)
	}
	m.UserComponents[user] = m.newUsercomponent()

	err = m.save()
	if err != nil {
		return err
	}

	return nil
}

// DeleteUsercomponent updates metadata deleting an user component
// XXX: this method isn't thread safe
func (m *Metadata) DeleteUsercomponent(user string) error {
	err := m.load()
	if err != nil {
		return err
	}

	_, ok := m.UserComponents[user]
	if !ok {
		return errors.Errorf("user component %s doesn't exist", user)
	}
	delete(m.UserComponents, user)
	err = m.save()
	if err != nil {
		return err
	}

	return nil
}

// ListUsercomponents returns an array with user compoment names for the deployment
func (m *Metadata) ListUsercomponents() ([]string, error) {
	err := m.load()
	if err != nil {
		return []string{}, err
	}

	return m.listUsercomponents()
}

// CheckUsercomponent checks if a user component is created
func (m *Metadata) CheckUsercomponent(user string) (bool, error) {
	err := m.load()
	if err != nil {
		return false, err
	}

	return m.checkUsercomponent(user)
}

// GetGlobalFlavour returns the Flavour attribute from Metadata
func (m *Metadata) GetGlobalFlavour() (string, error) {
	err := m.load()
	if err != nil {
		return "", err
	}

	return m.Flavour, nil
}

// SetGlobalFlavour saves the value of given Flavour attribute on Metadata
// XXX: this method isn't thread safe
func (m *Metadata) SetGlobalFlavour(flavour string) error {
	err := m.load()
	if err != nil {
		return err
	}

	// TODO: check if flavour is defined or its valid.
	m.Flavour = flavour

	err = m.save()
	if err != nil {
		return err
	}

	return nil
}

// GetUserFlavour returns the Flavour attribute from Metadata for an specified user.
func (m *Metadata) GetUserFlavour(user string) (string, error) {
	err := m.load()
	if err != nil {
		return "", err
	}

	ok, err := m.checkUsercomponent(user)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.Errorf("user component %s doesn't exist", user)
	}

	return m.UserComponents[user].Flavour, nil
}

// SetUserFlavour saves the value of given Flavour attribute on Metadata for an specified user.
// XXX: this method isn't thread safe
func (m *Metadata) SetUserFlavour(flavour string, user string) error {
	err := m.load()
	if err != nil {
		return err
	}

	ok, err := m.checkUsercomponent(user)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Errorf("user component %s doesn't exist", user)
	}

	userComponent := m.UserComponents[user]
	userComponent.Flavour = flavour
	m.UserComponents[user] = userComponent

	err = m.save()
	if err != nil {
		return err
	}

	return nil
}

func newMetadata(fs afero.Fs, varsPath string) *Metadata {
	return &Metadata{
		fs:       fs,
		filePath: filepath.Join(varsPath, metadataFileName),

		UserComponents: map[string]userComponent{},
		Plugins:        []globalPlugin{},
	}
}

func (m *Metadata) newUsercomponent() userComponent {
	return userComponent{
		Plugins: []userPlugin{},
		Flavour: "default",
	}
}

func (m *Metadata) load() error {
	data, err := afero.ReadFile(m.fs, m.filePath)
	if err != nil {
		return errors.Wrapf(err, "couldn't read file %s", m.filePath)
	}

	err = json.Unmarshal(data, m)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal json from file %s", m.filePath)
	}

	return nil
}

func (m *Metadata) save() error {
	data, err := json.MarshalIndent(*m, "", "  ")
	if err != nil {
		return errors.Wrap(err, "couldn't marshal json")
	}

	err = afero.WriteFile(m.fs, m.filePath, data, 0644)
	if err != nil {
		return errors.Wrapf(err, "couldn't write metadata file")
	}

	return nil
}

func (m *Metadata) listGlobalPlugins() ([]string, error) {
	list := []string{}
	for _, plugin := range m.Plugins {
		list = append(list, plugin.Name)
	}

	return list, nil
}

func (m *Metadata) listUserPlugins(user string) ([]string, error) {
	list := []string{}
	for _, plugin := range m.UserComponents[user].Plugins {
		list = append(list, plugin.Name)
	}

	return list, nil
}

func (m *Metadata) getGlobalPlugin(name string) (globalPlugin, error) {
	i, err := m.getGlobalPluginIndex(name)
	if err != nil {
		return globalPlugin{}, err
	}

	return m.Plugins[i], nil
}

func (m *Metadata) getGlobalPluginIndex(name string) (int, error) {
	for i, plugin := range m.Plugins {
		if plugin.Name == name {
			return i, nil
		}
	}

	return -1, errors.Errorf("global plugin %s doesn't exist", name)
}

func (m *Metadata) deleteGlobalPluginWithIndex(i int) {
	m.Plugins = append(m.Plugins[:i], m.Plugins[i+1:]...)
}

func (m *Metadata) globalPluginExists(name string) bool {
	for _, plugin := range m.Plugins {
		if plugin.Name == name {
			return true
		}
	}
	return false
}

func (m *Metadata) userPluginExists(name string, user string) bool {
	for _, plugin := range m.UserComponents[user].Plugins {
		if plugin.Name == name {
			return true
		}
	}
	return false
}

func (m *Metadata) checkUsercomponent(user string) (bool, error) {
	list, err := m.listUsercomponents()
	if err != nil {
		return false, err
	}

	_, ok := utils.FindString(list, user)

	return ok, nil
}

func (m *Metadata) listUsercomponents() ([]string, error) {
	keys := []string{}
	for k := range m.UserComponents {
		keys = append(keys, k)
	}

	return keys, nil
}

func (m *Metadata) getUserPlugin(name string, user string) (userPlugin, error) {
	i, err := m.getUserPluginIndex(name, user)
	if err != nil {
		return userPlugin{}, err
	}

	return m.UserComponents[user].Plugins[i], nil
}

func (m *Metadata) getUserPluginIndex(name string, user string) (int, error) {
	for i, plugin := range m.UserComponents[user].Plugins {
		if plugin.Name == name {
			return i, nil
		}
	}

	return -1, errors.Errorf("user plugin %s doesn't exist for user %s", name, user)
}

func (m *Metadata) deleteUserPluginWithIndex(i int, user string) {
	userComponent := m.UserComponents[user]
	userComponent.Plugins = append(userComponent.Plugins[:i], userComponent.Plugins[i+1:]...)
	m.UserComponents[user] = userComponent
}
