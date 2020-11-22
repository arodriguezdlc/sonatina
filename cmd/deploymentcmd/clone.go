package deploymentcmd

import (
	"fmt"

	"github.com/arodriguezdlc/sonatina/cmd/common"
	"github.com/arodriguezdlc/sonatina/manager"

	"github.com/spf13/cobra"
)

// CloneDeployment declares `sonatina clone deployment` command
var CloneDeployment = &cobra.Command{
	Use:   "deployment",
	Short: "Get a deployment from state repository",
	Args:  cobra.ExactArgs(1),
	RunE:  cloneDeploymentExecution,
}

func init() {
	CloneDeployment.Flags().StringVarP(&storageRepoURI, "storage-repo-uri", "s", "", "storage git repo uri")
	CloneDeployment.MarkFlagRequired("storage-repo-uri")
}

func cloneDeploymentExecution(command *cobra.Command, args []string) error {
	deployName := args[0]
	m := manager.GetManager()

	err := m.Clone(deployName, storageRepoURI)
	if err != nil {
		return err
	}

	err = common.SetCurrentDeployment(deployName)
	if err != nil {
		return err
	}

	fmt.Println("Cloned")
	return nil
}
