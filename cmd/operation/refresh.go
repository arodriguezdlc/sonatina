package operation

import (
	"fmt"

	"github.com/arodriguezdlc/sonatina/cmd/common"
	"github.com/arodriguezdlc/sonatina/manager"
	"github.com/spf13/cobra"
)

// Refresh declares `sonatina refresh` command
var Refresh = &cobra.Command{
	Use:   "refresh",
	Short: "Retrieve last changes from git repositories",
	RunE:  refreshExecution,
}

func init() {
	Refresh.Flags().StringVarP(&deployName, "deployment", "d", "", "deployment name")
}

func refreshExecution(command *cobra.Command, args []string) error {
	deployName, err := common.GetCurrentDeployment(deployName)
	if err != nil {
		return err
	}

	m := manager.GetManager()
	deploy, err := m.Get(deployName)
	if err != nil {
		return err
	}

	err = deploy.Pull()
	if err != nil {
		return err
	}

	fmt.Println("Refreshed")
	return nil
}
