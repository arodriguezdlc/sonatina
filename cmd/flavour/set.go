package flavour

import (
	"fmt"

	"github.com/arodriguezdlc/sonatina/cmd/common"
	"github.com/arodriguezdlc/sonatina/manager"
	"github.com/spf13/cobra"
)

// SetFlavour declares `sonatina set flavour` command
var SetFlavour = &cobra.Command{
	Use:   "flavour",
	Short: "",
	Long:  "",
	Args:  cobra.ExactArgs(1),
	RunE:  setFlavourExecution,
}

func init() {
	SetFlavour.Flags().StringVarP(&deployName, "deployment", "d", "", "deployment name")
	SetFlavour.Flags().StringVarP(&userComponent, "user-component", "c", "", "user component")
}

func setFlavourExecution(command *cobra.Command, args []string) error {
	flavour := args[0]

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
		err = deploy.SetFlavourGlobal(flavour)
		if err != nil {
			return err
		}
	} else {
		err = deploy.SetFlavourUser(flavour, userComponent)
		if err != nil {
			return err
		}
	}

	fmt.Println("Configured")
	return nil
}
