package operation

import (
	"os"
	"os/exec"

	"github.com/arodriguezdlc/sonatina/manager"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Edit declares `sonatina edit` command
var Edit = &cobra.Command{
	Use:   "edit",
	Short: "",
	Long:  "",
	RunE:  editExecution,
}

func init() {
	Edit.Flags().StringVarP(&deployName, "deployment", "d", "", "deployment name")
	Edit.MarkFlagRequired("deployment") // TODO: use current deployment by default and remove MarkFlagRequired

	Edit.Flags().StringVarP(&userComponent, "user-component", "c", "", "user component")
	Edit.Flags().StringVarP(&pluginName, "plugin", "p", "", "plugin")
}

func editExecution(command *cobra.Command, args []string) error {
	m := manager.GetManager()
	deploy, err := m.Get(deployName)
	if err != nil {
		return err
	}

	filepath, err := deploy.GetVariableFilepath("config", pluginName, userComponent)
	if err != nil {
		return err
	}

	err = openEditor(filepath)
	if err != nil {
		return err
	}

	return nil
}

func openEditor(filepath string) error {
	editor := viper.GetString("Editor")

	cmd := exec.Command(editor, filepath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "error executing the editor %s", editor)
	}

	return nil
}
