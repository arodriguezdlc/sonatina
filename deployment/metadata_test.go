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
		t.Fatal(err)
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
		t.Fatal(err)
	}

	expectedJSON := testMetadataReferenceJSON()

	obtainedJSONBytes, err := afero.ReadFile(fs, "/metadata.json")
	if err != nil {
		t.Fatal(err)
	}
	obtainedJSON := string(obtainedJSONBytes)

	if !reflect.DeepEqual(expectedJSON, obtainedJSON) {
		t.Errorf("Incorrect generated json metadata.\n\n Expected:\n%s\n\n Obtained:\n%s\n", expectedJSON, obtainedJSON)
	}
}

func TestListGlobalPlugins(t *testing.T) {
	fs := afero.NewMemMapFs()

	metadata := testNewMetadataWithData(fs)
	err := metadata.save()
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"plugin1", "plugin2"}
	obtained, err := metadata.ListGlobalPlugins()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, obtained) {
		t.Errorf("Incorrect global plugin list from metadata.\n\n Expected:\n%s\n\n Obtained:\n%s\n", expected, obtained)
	}
}

func TestListUserPlugins(t *testing.T) {
	fs := afero.NewMemMapFs()

	metadata := testNewMetadataWithData(fs)
	err := metadata.save()
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"plugin1", "plugin2"}
	obtained, err := metadata.ListUserPlugins("user1")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, obtained) {
		t.Errorf("Incorrect global plugin list from metadata.\n\n Expected:\n%s\n\n Obtained:\n%s\n", expected, obtained)
	}
}

func testNewMetadataEmpty(fs afero.Fs) Metadata {
	return Metadata{
		fs:       fs,
		filePath: "/metadata.json",
	}
}

func testNewMetadataWithData(fs afero.Fs) Metadata {
	metadata := Metadata{
		fs:       fs,
		filePath: "/metadata.json",

		TerraformVersion: "0.12.24",
		Repo:             "/code/repository",
		RepoPath:         "/", //TODO: support for specific path in repository
		Version:          "0.0.1",
		Commit:           "abcdefghijklmnopqrstuvwyz0123456789",
		Flavour:          "default",
		UserComponents:   map[string]userComponent{"user1": testNewUserComponent(), "user2": testNewUserComponent()},
		Plugins:          []globalPlugin{testNewGlobalPlugin("plugin1"), testNewGlobalPlugin("plugin2")},
	}

	return metadata
}

func testNewUserComponent() userComponent {
	userComponent := userComponent{
		Flavour: "default",
		Plugins: []userPlugin{testNewUserPlugin("plugin1"), testNewUserPlugin("plugin2")},
	}

	return userComponent
}

func testNewUserPlugin(name string) userPlugin {
	return userPlugin{
		Name: name,
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
	return afero.WriteFile(fs, "/metadata.json", []byte(testMetadataReferenceJSON()), 0644)
}

func testMetadataReferenceJSON() string {
	return `{
  "terraform_version": "0.12.24",
  "repo": "/code/repository",
  "repo_path": "/",
  "version": "0.0.1",
  "commit": "abcdefghijklmnopqrstuvwyz0123456789",
  "flavour": "default",
  "user_components": {
    "user1": {
      "plugins": [
        {
          "name": "plugin1"
        },
        {
          "name": "plugin2"
        }
      ],
      "flavour": "default"
    },
    "user2": {
      "plugins": [
        {
          "name": "plugin1"
        },
        {
          "name": "plugin2"
        }
      ],
      "flavour": "default"
    }
  },
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
