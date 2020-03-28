package deploymentcmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/arodriguezdlc/sonatina/manager"
	"github.com/spf13/cobra"
)

// ListDeployment declares `sonatina list deployments` command
var ListDeployment = &cobra.Command{
	Use:   "deployments",
	Short: "list deployments managed by sonatina",
	Long:  `TO DO`,
	Run:   listDeploymentExecution,
}

//To define flags
func init() {
}

func listDeploymentExecution(cmd *cobra.Command, args []string) {
	m := manager.GetManager()
	list, err := m.List()

	if err != nil {
		logrus.Fatalln(err)
	}

	fmt.Fprintln(os.Stdout, "DEPLOYMENTS:")
	for _, element := range list {
		fmt.Fprintln(os.Stdout, element)
	}
}
