package deployment

import (
	"encoding/json"
	"path/filepath"

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

func newMetadata(fs afero.Fs, varsPath string) *Metadata {
	return &Metadata{
		fs:       fs,
		filePath: filepath.Join(varsPath, metadataFileName),
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
