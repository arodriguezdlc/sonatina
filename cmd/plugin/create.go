package plugin

import (
	"errors"
	"fmt"

	"github.com/arodriguezdlc/sonatina/manager"

	"github.com/spf13/cobra"
)

// CreatePlugin declares `sonatina create plugin` command
var CreatePlugin = &cobra.Command{
	Use:   "plugin",
	Short: "create plugin",
	Long:  `create plugin`,
	Args:  cobra.ExactArgs(1),
	RunE:  createPluginExecution,
}

func init() {
	CreatePlugin.Flags().StringVarP(&deployName, "deployment", "d", "", "deployment name")
	CreatePlugin.MarkFlagRequired("deployment") // TODO: use current deployment by default and remove MarkFlagRequired

	CreatePlugin.Flags().StringVarP(&repoURI, "repo-uri", "r", "", "plugin git repo uri")
	CreatePlugin.Flags().StringVarP(&repoPath, "code-repo-path", "p", "", "code git repo path")

	CreatePlugin.Flags().StringVarP(&userComponent, "user-component", "c", "", "user component name")
}

func createPluginExecution(command *cobra.Command, args []string) error {
	pluginName := args[0]
	m := manager.GetManager()

	deploy, err := m.Get(deployName)
	if err != nil {
		return err
	}

	if userComponent == "" {
		if repoURI == "" { // Only required if it's a global plugin
			return errors.New("required flag(s) \"repo-uri\" not set")
		}
		err = deploy.CreatePluginGlobal(pluginName, repoURI, repoPath)
	} else {
		err = deploy.CreatePluginUser(pluginName, userComponent)
	}
	if err != nil {
		return err
	}

	fmt.Println("Created")
	return nil
}
