package deploymentcmd

import (
	"fmt"

	"github.com/arodriguezdlc/sonatina/manager"

	"github.com/spf13/cobra"
)

// DeleteDeployment declares `sonatina delete deployment` command
var DeleteDeployment = &cobra.Command{
	Use:   "deployment",
	Short: "Delete a specified deployment from local",
	Args:  cobra.ExactArgs(1),
	RunE:  deleteDeploymentExecution,
}

func deleteDeploymentExecution(command *cobra.Command, args []string) error {
	deployName := args[0]
	m := manager.GetManager()

	err := m.Delete(deployName)
	if err != nil {
		return err
	}

	fmt.Println("Deleted")
	return nil
}
