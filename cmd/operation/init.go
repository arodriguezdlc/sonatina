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
	Short: "Initialize a working directory",
	RunE:  initExecution,
}

func init() {
	Init.Flags().StringVarP(&deployName, "deployment", "d", "", "deployment name")
	Init.Flags().StringVarP(&userComponent, "user-component", "c", "", "user component")
}

func initExecution(command *cobra.Command, args []string) error {
	deployName, err := common.GetCurrentDeployment(deployName)
	if err != nil {
		return err
	}

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
	if userComponent == "" {
		err = init.RunGlobal()
	} else {
		err = init.RunUser(userComponent)
	}
	if err != nil {
		return err
	}

	return nil
}
