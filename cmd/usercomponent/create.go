package usercomponent

import (
	"fmt"

	"github.com/arodriguezdlc/sonatina/manager"

	"github.com/spf13/cobra"
)

// CreateUsercomponent declares `sonatina create usercomponent` command
var CreateUsercomponent = &cobra.Command{
	Use:   "usercomponent",
	Short: "create usercomponent",
	Long:  `create usercomponent`,
	Args:  cobra.ExactArgs(1),
	RunE:  createUsercomponentExecution,
}

func init() {
	CreateUsercomponent.Flags().StringVarP(&deployName, "deployment", "d", "", "deployment name")
	CreateUsercomponent.MarkFlagRequired("deployment") // TODO: use current deployment by default and remove MarkFlagRequired
}

func createUsercomponentExecution(command *cobra.Command, args []string) error {
	usercomponentName := args[0]
	m := manager.GetManager()

	deploy, err := m.Get(deployName)
	if err != nil {
		return err
	}

	err = deploy.CreateUsercomponent(usercomponentName)
	if err != nil {
		return err
	}

	fmt.Println("Created")
	return nil
}
