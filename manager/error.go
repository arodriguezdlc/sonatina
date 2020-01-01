package manager

import "fmt"

type DeploymentAlreadyExistsError struct {
	DeploymentName string
}

func (err DeploymentAlreadyExistsError) Error() string {
	return fmt.Sprintf("Deployment %v already exists", err.DeploymentName)
}
