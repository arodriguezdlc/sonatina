package workflow

import (
	"github.com/arodriguezdlc/sonatina/deployment"
	"github.com/arodriguezdlc/sonatina/terraformcli"
)

type DestroyWorkflow struct {
	Terraform  *terraformcli.Terraform
	Deployment deployment.Deployment
}

func Destroy(terraform *terraformcli.Terraform, deployment deployment.Deployment) *DestroyWorkflow {
	return &DestroyWorkflow{
		Terraform:  terraform,
		Deployment: deployment,
	}
}

func (i *DestroyWorkflow) RunGlobal(message string) error {
	executionPath, err := i.Deployment.GenerateWorkdirGlobal()
	if err != nil {
		return err
	}

	variableFiles, err := i.Deployment.GenerateVariablesGlobal()
	if err != nil {
		return err
	}

	stateFile := i.Deployment.StateFilePathGlobal()

	err = i.Terraform.Init(executionPath)
	if err != nil {
		return err
	}

	err = i.Terraform.Destroy(executionPath, variableFiles, stateFile)
	if err != nil {
		return err
	}

	err = i.Deployment.Push(message)
	if err != nil {
		return err
	}

	return nil
}

func (i *DestroyWorkflow) RunUser(message string, user string) error {
	executionPath, err := i.Deployment.GenerateWorkdirUser(user)
	if err != nil {
		return err
	}

	variableFiles, err := i.Deployment.GenerateVariablesUser(user)
	if err != nil {
		return err
	}

	stateFile := i.Deployment.StateFilePathUser(user)

	err = i.Terraform.Init(executionPath)
	if err != nil {
		return err
	}

	err = i.Terraform.Destroy(executionPath, variableFiles, stateFile)
	if err != nil {
		return err
	}

	err = i.Deployment.Push(message)
	if err != nil {
		return err
	}

	return nil
}
