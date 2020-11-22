package plugin

import (
	"fmt"
	"os"

	"github.com/arodriguezdlc/sonatina/cmd/common"
	"github.com/arodriguezdlc/sonatina/manager"

	"github.com/spf13/cobra"
)

// ListPlugins declares `sonatina list plugins` command
var ListPlugins = &cobra.Command{
	Use:   "plugins",
	Short: "List plugins added to deployment",
	RunE:  listPluginsExecution,
}

func init() {
	ListPlugins.Flags().StringVarP(&deployName, "deployment", "d", "", "deployment name")
	ListPlugins.Flags().StringVarP(&userComponent, "user-component", "c", "", "user component")
}

func listPluginsExecution(command *cobra.Command, args []string) error {
	var plugins []string

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
		plugins, err = deploy.ListPluginsGlobal()
		if err != nil {
			return err
		}
	} else {
		plugins, err = deploy.ListPluginsUser(userComponent)
		if err != nil {
			return err
		}
	}

	fmt.Fprintln(os.Stdout, "PLUGINS:")
	for _, element := range plugins {
		_, err = fmt.Fprintln(os.Stdout, element)
		if err != nil {
			return err
		}
	}

	return nil
}
