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

func (i *DestroyWorkflow) RunGlobal() error {
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

	return nil
}

// TODO
func (i *DestroyWorkflow) RunUser(user string) error {
	return nil
}
