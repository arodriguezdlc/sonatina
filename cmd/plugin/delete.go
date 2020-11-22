package plugin

import (
	"fmt"

	"github.com/arodriguezdlc/sonatina/cmd/common"
	"github.com/arodriguezdlc/sonatina/manager"

	"github.com/spf13/cobra"
)

// DeletePlugin declares `sonatina delete plugin` command
var DeletePlugin = &cobra.Command{
	Use:   "plugin",
	Short: "Remove a specified plugin from deployment",
	Args:  cobra.ExactArgs(1),
	RunE:  deletePluginExecution,
}

func init() {
	DeletePlugin.Flags().StringVarP(&deployName, "deployment", "d", "", "deployment name")
	DeletePlugin.Flags().StringVarP(&userComponent, "user-component", "c", "", "user component name")
}

func deletePluginExecution(command *cobra.Command, args []string) error {
	pluginName := args[0]

	deployName, err := common.GetCurrentDeployment(deployName)
	if err != nil {
		return err
	}

	m := manager.GetManager()
	deploy, err := m.Get(deployName)
	if err != nil {
		return err
	}

	if userComponent == "" {
		err = deploy.DeletePluginGlobal(pluginName)
	} else {
		err = deploy.DeletePluginUser(pluginName, userComponent)
	}
	if err != nil {
		return err
	}

	fmt.Println("Deleted")
	return nil
}
