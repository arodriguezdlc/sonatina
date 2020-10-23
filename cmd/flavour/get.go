package flavour

import (
	"fmt"

	"github.com/arodriguezdlc/sonatina/cmd/common"
	"github.com/arodriguezdlc/sonatina/manager"
	"github.com/spf13/cobra"
)

// GetFlavour declares `sonatina get flavour` command
var GetFlavour = &cobra.Command{
	Use:   "flavour",
	Short: "",
	Long:  "",
	RunE:  getFlavourExecution,
}

func init() {
	GetFlavour.Flags().StringVarP(&deployName, "deployment", "d", "", "deployment name")
	GetFlavour.Flags().StringVarP(&userComponent, "user-component", "c", "", "user component")
}

func getFlavourExecution(command *cobra.Command, args []string) error {
	var flavour string

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
		flavour, err = deploy.GetFlavourGlobal()
		if err != nil {
			return err
		}
	} else {
		flavour, err = deploy.GetFlavourUser(userComponent)
		if err != nil {
			return err
		}
	}

	fmt.Println(flavour)
	return nil
}
