package manager

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/arodriguezdlc/sonatina/deployment"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

type ManagerInterface interface {
	List() ([]string, error)
	Get(name string) (deployment.Deployment, error)
	Add(name string, storageRepoURI string, codeRepoURI string) (deployment.Deployment, error)
	Delete(name string) error
}

var manager ManagerInterface

// Uses manager from configuration
func init() {
}

// InitializeManager creates a deployment manager object and saves it
// in the manager global variable
func InitializeManager(fs afero.Fs) error {
	connector := viper.GetString("ManagerConnector")
	var err error
	switch connector {
	case "yaml":
		log.Infoln("Initialize YAML based Manager")
		manager, err = NewManagerYaml(fs, viper.GetString("DeploymentsPath"), viper.GetString("DeploymentsFilename"))
		return err
	default:
		manager = nil
		return errors.New("Unsupported manager connector [" + connector + "]")
	}
}

// GetManager returns the manager global variable, that was initialized by
// InitializeManager function
func GetManager() ManagerInterface {
	return manager
}
