package deployment

import (
	"encoding/json"
	"path/filepath"

	"github.com/arodriguezdlc/sonatina/utils"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

const metadataFileName string = "metadata.json"

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

func (m *Metadata) ListGlobalPlugins() []string {
	list := []string{}
	for _, plugin := range m.Plugins {
		list = append(list, plugin.Name)
	}
	return list
}

func (m *Metadata) ListUserPlugins(user string) []string {
	list := []string{}
	for _, plugin := range m.UserComponents[user].Plugins {
		list = append(list, plugin.Name)
	}
	return list
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
	keys := []string{}

	err := m.load()
	if err != nil {
		return keys, err
	}

	for k := range m.UserComponents {
		keys = append(keys, k)
	}

	return keys, nil
}

// CheckUsercomponent checks if a user component is created
func (m *Metadata) CheckUsercomponent(user string) (bool, error) {
	list, err := m.ListUsercomponents()
	if err != nil {
		return false, err
	}

	_, ok := utils.FindString(list, user)

	return ok, nil
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
