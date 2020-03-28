package deployment

import (
	"encoding/json"

	"github.com/spf13/afero"
)

type metadata struct {
	fs       afero.Fs
	filePath string

	Name           string          `json:"name"`
	Repo           string          `json:"repo"`
	RepoPath       string          `json:"repo_path"`
	Version        string          `json:"version"`
	Commit         string          `json:"commit"`
	Flavour        string          `json:"flavour"`
	UserComponents []userComponent `json:"user_components"`
	Plugins        []globalPlugin  `json:"plugins"`
}

type userComponent struct {
	Name    string       `json:"name"`
	Plugins []userPlugin `json:"plugins"`
}

type globalPlugin struct {
	Name     string `json:"name"`
	Repo     string `json:"repo"`
	RepoPath string `json:"repo_path"`
	Version  string `json:"version"`
	Commit   string `json:"commit"`
}

type userPlugin struct {
	Name    string `json:"name"`
	Flavour string `json:"flavour"`
}

func newMetadata(deployment DeploymentImpl) metadata {
	metadata := metadata{
		fs:       deployment.fs,
		filePath: deployment.Vars.path + "/metadata.yml",

		Name:           deployment.Name,
		Repo:           deployment.CodeRepoURL,
		RepoPath:       "/", //TODO: support for specific path in repository
		Version:        "",
		Commit:         "",
		Flavour:        "",
		UserComponents: nil,
		Plugins:        nil,
	}

	return metadata
}

func (m *metadata) load() error {
	data, err := afero.ReadFile(m.fs, m.filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, m)
	if err != nil {
		return err
	}

	return nil
}

func (m *metadata) save() error {
	data, err := json.MarshalIndent(*m, "", "  ")
	if err != nil {
		return err
	}

	err = afero.WriteFile(m.fs, m.filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
