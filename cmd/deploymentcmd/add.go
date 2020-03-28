package deploymentcmd

import (
	"fmt"

	"github.com/arodriguezdlc/sonatina/manager"
	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// AddDeployment declares `sonatina add deployment` command
var AddDeployment = &cobra.Command{
	Use:   "deployment",
	Short: "add deployment",
	Long:  `add deployment`,
	Args:  cobra.ExactArgs(1),
	Run:   addDeploymentExecution,
}

//To define flags
var storageRepoURI string
var codeRepoURI string

func init() {
	AddDeployment.Flags().StringVarP(&storageRepoURI, "storage-repo-uri", "s", "", "storage git repo uri")
	AddDeployment.MarkFlagRequired("storage-repo-uri")

	AddDeployment.Flags().StringVarP(&codeRepoURI, "code-repo-uri", "c", "", "code git repo uri")
	AddDeployment.MarkFlagRequired("code-repo-uri")
}

func addDeploymentExecution(cmd *cobra.Command, args []string) {
	var err error

	deployName := args[0]
	m := manager.GetManager()

	if _, err = m.Add(deployName, storageRepoURI, codeRepoURI); err != nil {
		logrus.Fatalln(err)
	}

	fmt.Println("Created")
}
