package manager

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

func TestInitializeManagerYaml(t *testing.T) {
	fs := afero.NewMemMapFs()
	viper.Set("ManagerConnector", "json")

	err := InitializeManager(fs)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInitializeManagerIncorrect(t *testing.T) {
	fs := afero.NewMemMapFs()
	viper.Set("ManagerConnector", "Invalid")

	err := InitializeManager(fs)
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
