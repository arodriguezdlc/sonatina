package workflow

import (
	"github.com/arodriguezdlc/sonatina/deployment"
	"github.com/arodriguezdlc/sonatina/terraformcli"
)

type InitWorkflow struct {
	Terraform  *terraformcli.Terraform
	Deployment deployment.Deployment
}

func Init(terraform *terraformcli.Terraform, deployment deployment.Deployment) *InitWorkflow {
	return &InitWorkflow{
		Terraform:  terraform,
		Deployment: deployment,
	}
}

func (i *InitWorkflow) RunGlobal() error {
	err := i.Deployment.Pull()
	if err != nil {
		return err
	}

	executionPath, err := i.Deployment.GenerateWorkdirGlobal()
	if err != nil {
		return err
	}

	_, err = i.Deployment.GenerateVariablesGlobal()
	if err != nil {
		return err
	}

	err = i.Terraform.Init(executionPath)
	if err != nil {
		return err
	}

	return nil
}

func (i *InitWorkflow) RunUser(user string) error {
	err := i.Deployment.Pull()
	if err != nil {
		return err
	}

	executionPath, err := i.Deployment.GenerateWorkdirUser(user)
	if err != nil {
		return err
	}

	_, err = i.Deployment.GenerateVariablesUser(user)
	if err != nil {
		return err
	}

	err = i.Terraform.Init(executionPath)
	if err != nil {
		return err
	}

	return nil
}
