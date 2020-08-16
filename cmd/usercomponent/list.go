package usercomponent

import (
	"fmt"
	"os"

	"github.com/arodriguezdlc/sonatina/cmd/common"
	"github.com/arodriguezdlc/sonatina/manager"

	"github.com/spf13/cobra"
)

// ListUsercomponents declares `sonatina list usercomponent` command
var ListUsercomponents = &cobra.Command{
	Use:   "usercomponents",
	Short: "list usercomponents",
	Long:  `list usercomponents`,
	RunE:  listUsercomponentsExecution,
}

func init() {
	ListUsercomponents.Flags().StringVarP(&deployName, "deployment", "d", "", "deployment name")
}

func listUsercomponentsExecution(command *cobra.Command, args []string) error {
	deployName, err := common.GetCurrentDeployment(deployName)
	if err != nil {
		return err
	}

	m := manager.GetManager()
	deploy, err := m.Get(deployName)
	if err != nil {
		return err
	}

	list, err := deploy.ListUsercomponents()
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, "USER COMPONENTS:")
	for _, element := range list {
		_, err = fmt.Fprintln(os.Stdout, element)
		if err != nil {
			return err
		}
	}

	return nil
}
