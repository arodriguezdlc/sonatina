package manager

import (
	"testing"

	"github.com/spf13/afero"
)

func TestInitializeManagerYaml(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := InitializeManager(fs, "json")
	if err != nil {
		t.Fatal(err)
	}
}

func TestInitializeManagerIncorrect(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := InitializeManager(fs, "invalid")
	if _, ok := err.(ManagerUnsupportedConnectorError); !ok {
		t.Errorf("Expected ManagerUnsupportedConnectorError, obtained %v", err)
	}
}

func TestGetManager(t *testing.T) {
	// TODO
}

func TestList(t *testing.T) {
	// TODO
}
