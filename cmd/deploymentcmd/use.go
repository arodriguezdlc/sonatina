package deploymentcmd

import (
	"fmt"

	"github.com/arodriguezdlc/sonatina/cmd/common"
	"github.com/arodriguezdlc/sonatina/manager"

	"github.com/spf13/cobra"
)

// UseDeployment declares `sonatina use deployment` command
var UseDeployment = &cobra.Command{
	Use:   "deployment",
	Short: "use deployment",
	Long:  `use deployment`,
	Args:  cobra.ExactArgs(1),
	RunE:  useDeploymentExecution,
}

func useDeploymentExecution(command *cobra.Command, args []string) error {
	deployName := args[0]
	m := manager.GetManager()

	_, err := m.Get(deployName)
	if err != nil {
		return err
	}

	err = common.SetCurrentDeployment(deployName)
	if err != nil {
		return err
	}

	fmt.Println("Configured")
	return nil
}
