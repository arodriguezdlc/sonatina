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
	executionPath, err := i.Deployment.GenerateWorkdirGlobal()
	if err != nil {
		return err
	}

	err = i.Terraform.Init(executionPath)
	if err != nil {
		return err
	}

	return nil
}

// TODO
func (i *InitWorkflow) RunUser(user string) error {
	return nil
}

func (i *InitWorkflow) configureCommand() {

}
