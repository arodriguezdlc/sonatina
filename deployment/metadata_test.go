package deployment

import (
	"reflect"
	"testing"

	"github.com/spf13/afero"
)

func TestLoadMetadata(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := testWriteMetadataReferenceFile(fs)
	if err != nil {
		t.Fatal(err)
	}

	expectedMetadata := testNewMetadataWithData(fs)
	loadedMetadata := testNewMetadataEmpty(fs)

	err = loadedMetadata.load()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expectedMetadata, loadedMetadata) {
		t.Errorf("Incorrect loaded metadata.\n\n Expected: %v\n\n Obtained: %v\n", expectedMetadata, loadedMetadata)
	}
}

func TestSaveMetadata(t *testing.T) {
	fs := afero.NewMemMapFs()

	metadata := testNewMetadataWithData(fs)

	err := metadata.save()
	if err != nil {
		t.Error(err)
	}

	expectedJSON := testMetadataReferenceJSON()

	obtainedJSONBytes, err := afero.ReadFile(fs, "/metadata.yml")
	if err != nil {
		t.Error(err)
	}
	obtainedJSON := string(obtainedJSONBytes)

	if !reflect.DeepEqual(expectedJSON, obtainedJSON) {
		t.Errorf("Incorrect generated json metadata.\n\n Expected:\n%s\n\n Obtained:\n%s\n", expectedJSON, obtainedJSON)
	}
}

func testNewMetadataEmpty(fs afero.Fs) metadata {
	return metadata{
		fs:       fs,
		filePath: "/metadata.yml",
	}
}

func testNewMetadataWithData(fs afero.Fs) metadata {
	metadata := metadata{
		fs:       fs,
		filePath: "/metadata.yml",

		Name:           "test",
		Repo:           "/code/repository",
		RepoPath:       "/", //TODO: support for specific path in repository
		Version:        "0.0.1",
		Commit:         "abcdefghijklmnopqrstuvwyz0123456789",
		Flavour:        "default",
		UserComponents: []userComponent{testNewUserComponent("user1"), testNewUserComponent("user2")},
		Plugins:        []globalPlugin{testNewGlobalPlugin("plugin1"), testNewGlobalPlugin("plugin2")},
	}

	return metadata
}

func testNewUserComponent(name string) userComponent {
	userComponent := userComponent{
		Name:    name,
		Plugins: []userPlugin{testNewUserPlugin("plugin1"), testNewUserPlugin("plugin2")},
	}

	return userComponent
}

func testNewUserPlugin(name string) userPlugin {
	return userPlugin{
		Name:    name,
		Flavour: "default",
	}
}

func testNewGlobalPlugin(name string) globalPlugin {
	return globalPlugin{
		Name:     name,
		Repo:     "/repo/" + name,
		RepoPath: "/",
		Version:  "0.0.1",
		Commit:   "abcdefghijklmnopqrstuvwyz0123456789",
	}
}

func testWriteMetadataReferenceFile(fs afero.Fs) error {
	return afero.WriteFile(fs, "/metadata.yml", []byte(testMetadataReferenceJSON()), 0644)
}

func testMetadataReferenceJSON() string {
	return `{
  "name": "test",
  "repo": "/code/repository",
  "repo_path": "/",
  "version": "0.0.1",
  "commit": "abcdefghijklmnopqrstuvwyz0123456789",
  "flavour": "default",
  "user_components": [
    {
      "name": "user1",
      "plugins": [
        {
          "name": "plugin1",
          "flavour": "default"
        },
        {
          "name": "plugin2",
          "flavour": "default"
        }
      ]
    },
    {
      "name": "user2",
      "plugins": [
        {
          "name": "plugin1",
          "flavour": "default"
        },
        {
          "name": "plugin2",
          "flavour": "default"
        }
      ]
    }
  ],
  "plugins": [
    {
      "name": "plugin1",
      "repo": "/repo/plugin1",
      "repo_path": "/",
      "version": "0.0.1",
      "commit": "abcdefghijklmnopqrstuvwyz0123456789"
    },
    {
      "name": "plugin2",
      "repo": "/repo/plugin2",
      "repo_path": "/",
      "version": "0.0.1",
      "commit": "abcdefghijklmnopqrstuvwyz0123456789"
    }
  ]
}`
}
