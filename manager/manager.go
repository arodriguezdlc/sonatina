package manager

import (
	"github.com/arodriguezdlc/sonatina/deployment"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

// Manager interface that provides methods to manage differents
// deployments.
type Manager interface {
	List() ([]string, error)
	Get(name string) (deployment.Deployment, error)
	Create(name string, storageRepoURI string, codeRepoURI string, codeRepoPath string,
		terraformVersion string, flavour string) error
	Clone(name string, storageRepoURI string) error
	Delete(name string) error
}

// Single manager on sonatina execution. It's initialized at program start,
// then can be retrieved using GetManager method.
var manager Manager

// InitializeManager creates a deployment manager object and saves it
// in the manager global variable
func InitializeManager(fs afero.Fs, connector string) error {
	var err error
	switch connector {
	case "json":
		logrus.Infoln("Initialize JSON based Manager")
		manager, err = newManagerJSON(fs, viper.GetString("DeploymentsPath"), viper.GetString("DeploymentsFilename"))
		return err
	default:
		manager = nil
		return ManagerUnsupportedConnectorError{connector}
	}
}

// GetManager returns the manager global variable, that was initialized by
// InitializeManager function
func GetManager() Manager {
	return manager
}
