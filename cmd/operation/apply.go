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
	Args:  cobra.ExactArgs(1),
	RunE:  applyExecution,
}

func init() {
	Apply.Flags().StringVarP(&deployName, "deployment", "d", "", "deployment name")
	Apply.Flags().BoolVarP(&pull, "pull", "p", false, "enable pull before apply")
	Apply.Flags().StringVarP(&userComponent, "user-component", "c", "", "user component")
}

func applyExecution(command *cobra.Command, args []string) error {
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

	if pull {
		err = deploy.Pull()
		if err != nil {
			return err
		}
	}

	apply := workflow.Apply(terraform, deploy)
	if userComponent == "" { // TODO: check if is a valid user Component
		err = apply.RunGlobal(message)
	} else {
		err = apply.RunUser(message, userComponent)
	}
	if err != nil {
		return err
	}

	return nil
}
