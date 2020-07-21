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
	RunE:  destroyExecution,
}

func init() {
	Destroy.Flags().StringVarP(&deployName, "deployment", "d", "", "deployment name")
	Destroy.MarkFlagRequired("deployment") // TODO: use current deployment by default and remove MarkFlagRequired

	Destroy.Flags().StringVarP(&userComponent, "user-component", "c", "", "user component")
}

func destroyExecution(command *cobra.Command, args []string) error {
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
		err = destroy.RunGlobal()
	} else {
		err = destroy.RunUser(userComponent)
	}
	if err != nil {
		return err
	}

	return nil
}
