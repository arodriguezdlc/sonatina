package deploymentcmd

import (
	"fmt"
	"os"

	"github.com/arodriguezdlc/sonatina/cmd/common"
	"github.com/arodriguezdlc/sonatina/manager"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ListDeployment declares `sonatina list deployments` command
var ListDeployment = &cobra.Command{
	Use:   "deployments",
	Short: "list deployments managed by sonatina",
	Long:  `TO DO`,
	RunE:  listDeploymentExecution,
}

func listDeploymentExecution(command *cobra.Command, args []string) error {
	m := manager.GetManager()

	list, err := m.List()
	if err != nil {
		return err
	}

	current, err := common.GetCurrentDeployment("")
	if err != nil {
		logrus.WithError(err).Warning("couldn't get current deployment")
	}

	fmt.Fprintln(os.Stdout, "DEPLOYMENTS:")
	for _, element := range list {
		prefix := " - "
		if element == current {
			prefix = " * "
		}
		_, err = fmt.Fprintln(os.Stdout, prefix+element)
		if err != nil {
			return err
		}
	}

	return nil
}
