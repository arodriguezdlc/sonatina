package operation

import (
	"github.com/arodriguezdlc/sonatina/cmd/common"
	"github.com/arodriguezdlc/sonatina/manager"
	"github.com/arodriguezdlc/sonatina/workflow"
	"github.com/spf13/cobra"
)

// Destroy declares `sonatina destroy` command
var Destroy = &cobra.Command{
	Use:   "destroy",
	Short: "",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	RunE:  destroyExecution,
}

func init() {
	Destroy.Flags().StringVarP(&deployName, "deployment", "d", "", "deployment name")
	Destroy.Flags().StringVarP(&userComponent, "user-component", "c", "", "user component")
}

func destroyExecution(command *cobra.Command, args []string) error {
	message := args[0]

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

	destroy := workflow.Destroy(terraform, deploy)
	if userComponent == "" { // TODO: check if is a valid user Component
		err = destroy.RunGlobal(message)
	} else {
		err = destroy.RunUser(message, userComponent)
	}
	if err != nil {
		return err
	}

	return nil
}
