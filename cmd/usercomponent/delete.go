package usercomponent

import (
	"fmt"

	"github.com/arodriguezdlc/sonatina/cmd/common"
	"github.com/arodriguezdlc/sonatina/manager"

	"github.com/spf13/cobra"
)

// CreateUsercomponent declares `sonatina create usercomponent` command
var DeleteUsercomponent = &cobra.Command{
	Use:   "usercomponent",
	Short: "delete usercomponent",
	Long:  `delete usercomponent`,
	Args:  cobra.ExactArgs(1),
	RunE:  deleteUsercomponentExecution,
}

func init() {
	DeleteUsercomponent.Flags().StringVarP(&deployName, "deployment", "d", "", "deployment name")
}

func deleteUsercomponentExecution(command *cobra.Command, args []string) error {
	usercomponentName := args[0]

	deployName, err := common.GetCurrentDeployment(deployName)
	if err != nil {
		return err
	}

	m := manager.GetManager()
	deploy, err := m.Get(deployName)
	if err != nil {
		return err
	}

	err = deploy.DeleteUsercomponent(usercomponentName)
	if err != nil {
		return err
	}

	fmt.Println("Deleted")
	return nil
}
