package manager

import "fmt"

type ManagerUnsupportedConnectorError struct {
	Connector string
}

func (err ManagerUnsupportedConnectorError) Error() string {
	return fmt.Sprintf("Connector %v is not supported", err.Connector)
}

type DeploymentAlreadyExistsError struct {
	DeploymentName string
}

func (err DeploymentAlreadyExistsError) Error() string {
	return fmt.Sprintf("Deployment %v already exists", err.DeploymentName)
}

type DeploymentDoNotExistsError struct {
	DeploymentName string
}

func (err DeploymentDoNotExistsError) Error() string {
	return fmt.Sprintf("Deployment %v doesn't exist", err.DeploymentName)
}
