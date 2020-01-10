package deployment

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type Deployment interface {
}

type DeploymentImpl struct {
	Name  string
	State State
	Vars  Vars
	/*variables (vtd)
	state
	htd
	workingDirectory*/
}

// NewDeploymentImpl creates and initializes a new Deployment object
func NewDeploymentImpl(name string, storageRepoURL string, codeRepoURL string, fs afero.Fs, deploymentPath string) (Deployment, error) {
	var err error
	var vars Vars
	var state State

	//Create deployment directory
	if err = fs.MkdirAll(deploymentPath, 0700); err != nil {
		return nil, err
	}

	if vars, err = NewVarsGit(fs, deploymentPath+"/variables", storageRepoURL); err != nil {
		rollbackNewDeployment(fs, deploymentPath)
		return nil, err
	}

	if state, err = NewStateGit(fs, deploymentPath+"/state", storageRepoURL); err != nil {
		rollbackNewDeployment(fs, deploymentPath)
		return nil, err
	}

	return DeploymentImpl{
		Name:  name,
		Vars:  vars,
		State: state,
	}, nil
}

func rollbackNewDeployment(fs afero.Fs, deploymentPath string) {
	logrus.Debugln("rollbackNewDeployment deleting " + deploymentPath)
	if err := fs.RemoveAll(deploymentPath); err != nil {
		logrus.Errorln("Can't execute rollback for NewDeploymentImplementation") //TODO: improve error message
	}
}
