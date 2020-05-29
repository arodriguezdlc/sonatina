package operation

import (
	"github.com/arodriguezdlc/sonatina/cmd/common"
	"github.com/arodriguezdlc/sonatina/manager"
	"github.com/arodriguezdlc/sonatina/workflow"
	"github.com/spf13/cobra"
)

// Apply declares `sonatina apply` command
var Apply = &cobra.Command{
	Use:   "apply",
	Short: "",
	Long:  "",
	RunE:  applyExecution,
}

func init() {
	Apply.Flags().StringVarP(&deployName, "deployment", "d", "", "deployment name")
	Apply.MarkFlagRequired("deployment") // TODO: use current deployment by default and remove MarkFlagRequired

	Apply.Flags().StringVarP(&userComponent, "user-component", "c", "", "user component")
}

func applyExecution(command *cobra.Command, args []string) error {
	m := manager.GetManager()
	deploy, err := m.Get(deployName)
	if err != nil {
		return err
	}

	terraform, err := common.InitializeTerraform(deploy)
	if err != nil {
		return err
	}

	apply := workflow.Apply(terraform, deploy)
	if userComponent == "" { // TODO: check if is a valid user Component
		err = apply.RunGlobal()
	} else {
		err = apply.RunUser(userComponent)
	}
	if err != nil {
		return err
	}

	return nil
}
