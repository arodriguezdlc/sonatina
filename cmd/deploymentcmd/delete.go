package deploymentcmd

import (
	"fmt"

	"github.com/arodriguezdlc/sonatina/manager"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// DeleteDeployment declares `sonatina delete deployment` command
var DeleteDeployment = &cobra.Command{
	Use:   "deployment",
	Short: "delete deployment",
	Long:  `delete deployment`,
	Args:  cobra.ExactArgs(1),
	Run:   deleteDeploymentExecution,
}

func deleteDeploymentExecution(cmd *cobra.Command, args []string) {
	var err error

	deployName := args[0]
	m := manager.GetManager()

	if err = m.Delete(deployName); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Deleted")
}
