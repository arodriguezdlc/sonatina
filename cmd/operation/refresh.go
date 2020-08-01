package operation

import (
	"fmt"

	"github.com/arodriguezdlc/sonatina/manager"
	"github.com/spf13/cobra"
)

// Refresh declares `sonatina refresh` command
var Refresh = &cobra.Command{
	Use:   "refresh",
	Short: "",
	Long:  "",
	RunE:  refreshExecution,
}

func init() {
	Refresh.Flags().StringVarP(&deployName, "deployment", "d", "", "deployment name")
	Refresh.MarkFlagRequired("deployment") // TODO: use current deployment by default and remove MarkFlagRequired
}

func refreshExecution(command *cobra.Command, args []string) error {

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
