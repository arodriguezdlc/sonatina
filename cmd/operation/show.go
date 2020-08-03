package operation

import (
	"fmt"

	"github.com/arodriguezdlc/sonatina/manager"
	"github.com/spf13/cobra"
)

// Show declares `sonatina show` command
var Show = &cobra.Command{
	Use:   "show",
	Short: "",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	RunE:  showExecution,
}

func init() {
	Show.Flags().StringVarP(&deployName, "deployment", "d", "", "deployment name")
	Show.MarkFlagRequired("deployment") // TODO: use current deployment by default and remove MarkFlagRequired

	Show.Flags().StringVarP(&userComponent, "user-component", "c", "", "user component")
	Show.Flags().StringVarP(&pluginName, "plugin", "p", "", "plugin")
}

func showExecution(command *cobra.Command, args []string) error {
	kind := args[0]

	m := manager.GetManager()
	deploy, err := m.Get(deployName)
	if err != nil {
		return err
	}

	content, err := deploy.ReadVariableFilepath(kind, pluginName, userComponent)
	if err != nil {
		return err
	}

	fmt.Print(content)
	return nil
}
