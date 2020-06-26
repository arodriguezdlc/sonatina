package deploymentcmd

import (
	"fmt"
	"os"

	"github.com/arodriguezdlc/sonatina/manager"
	"github.com/spf13/cobra"
)

// ListDeployment declares `sonatina list deployments` command
var ListDeployment = &cobra.Command{
	Use:   "deployments",
	Short: "list deployments managed by sonatina",
	Long:  `TO DO`,
	RunE:  listDeploymentExecution,
}

func listDeploymentExecution(command *cobra.Command, args []string) error {
	m := manager.GetManager()

	list, err := m.List()
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, "DEPLOYMENTS:")
	for _, element := range list {
		_, err = fmt.Fprintln(os.Stdout, element)
		if err != nil {
			return err
		}
	}

	return nil
}
