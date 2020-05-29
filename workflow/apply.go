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

func (i *ApplyWorkflow) RunGlobal() error {
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

	return nil
}

// TODO
func (i *ApplyWorkflow) RunUser(user string) error {
	return nil
}
