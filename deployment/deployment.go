package deployment

import (
	"github.com/spf13/afero"
)

// Deployment representation
type Deployment interface {
	//Add()
	//List()
	//Delete()
}

// DeploymentImpl implements Deployment interface
type DeploymentImpl struct {
	Name  string
	State State
	Vars  Vars
	/*variables (vtd)
	state
	htd
	workingDirectory*/
}

// NewDeploymentImpl creates and initializes a new DeploymentImpl object
func NewDeploymentImpl(name string, repoURL string, fs afero.Fs, deploymentPath string) (Deployment, error) {
	var err error
	var vars Vars
	var state State

	if vars, err = NewVarsGit(fs, deploymentPath+"/variables", repoURL); err != nil {
		return nil, err
	}

	if state, err = NewStateGit(fs, deploymentPath+"/state", repoURL); err != nil {
		return nil, err
	}

	return DeploymentImpl{
		Name:  name,
		Vars:  vars,
		State: state,
	}, nil
}
