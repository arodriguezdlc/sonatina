package workflow

import (
	"github.com/arodriguezdlc/sonatina/deployment"
	"github.com/arodriguezdlc/sonatina/terraformcli"
)

type ApplyWorkflow struct {
	Terraform  *terraformcli.Terraform
	Deployment deployment.Deployment
}

func Apply(terraform *terraformcli.Terraform, deployment deployment.Deployment) *ApplyWorkflow {
	return &ApplyWorkflow{
		Terraform:  terraform,
		Deployment: deployment,
	}
}

func (i *ApplyWorkflow) RunGlobal(message string) error {
	err := i.Deployment.Pull()
	if err != nil {
		return err
	}

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

	err = i.Terraform.Apply(executionPath, variableFiles, stateFile)
	if err != nil {
		return err
	}

	err = i.Deployment.Push(message)
	if err != nil {
		return err
	}

	return nil
}

func (i *ApplyWorkflow) RunUser(message string, user string) error {
	err := i.Deployment.Pull()
	if err != nil {
		return err
	}

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

	err = i.Terraform.Apply(executionPath, variableFiles, stateFile)
	if err != nil {
		return err
	}

	err = i.Deployment.Push(message)
	if err != nil {
		return err
	}

	return nil
}
