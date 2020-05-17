package operation

import (
	"github.com/arodriguezdlc/sonatina/cmd/common"
	"github.com/arodriguezdlc/sonatina/manager"
	"github.com/arodriguezdlc/sonatina/workflow"
	"github.com/spf13/cobra"
)

// Init declares `sonatina init` command
var Init = &cobra.Command{
	Use:   "init",
	Short: "",
	Long:  "",
	RunE:  initExecution,
}

//To define flags
var deployName string
var userComponent string

func init() {
	Init.Flags().StringVarP(&deployName, "deployment", "d", "", "deployment name")
	Init.MarkFlagRequired("deployment") // TODO: use current deployment by default and remove MarkFlagRequired

	Init.Flags().StringVarP(&userComponent, "user-component", "c", "", "user component")
}

func initExecution(command *cobra.Command, args []string) error {
	m := manager.GetManager()
	deploy, err := m.Get(deployName)
	if err != nil {
		return err
	}

	terraform, err := common.InitializeTerraform(deploy)
	if err != nil {
		return err
	}

	init := workflow.Init(terraform, deploy)
	if userComponent == "" { // TODO: check if is a valid user Component
		err = init.RunGlobal()
	} else {
		err = init.RunUser(userComponent)
	}
	if err != nil {
		return err
	}

	return nil
}
