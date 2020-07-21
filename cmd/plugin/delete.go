package plugin

import (
	"fmt"

	"github.com/arodriguezdlc/sonatina/manager"

	"github.com/spf13/cobra"
)

// DeletePlugin declares `sonatina delete plugin` command
var DeletePlugin = &cobra.Command{
	Use:   "plugin",
	Short: "delete plugin",
	Long:  `delete plugin`,
	Args:  cobra.ExactArgs(1),
	RunE:  deletePluginExecution,
}

func init() {
	DeletePlugin.Flags().StringVarP(&deployName, "deployment", "d", "", "deployment name")
	DeletePlugin.MarkFlagRequired("deployment") // TODO: use current deployment by default and remove MarkFlagRequired

	DeletePlugin.Flags().StringVarP(&userComponent, "user-component", "c", "", "user component name")
}

func deletePluginExecution(command *cobra.Command, args []string) error {
	pluginName := args[0]
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
