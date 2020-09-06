package deploymentcmd

import (
	"fmt"

	"github.com/arodriguezdlc/sonatina/cmd/common"
	"github.com/arodriguezdlc/sonatina/manager"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// CreateDeployment declares `sonatina create deployment` command
var CreateDeployment = &cobra.Command{
	Use:   "deployment",
	Short: "create deployment",
	Long:  `create deployment`,
	Args:  cobra.ExactArgs(1),
	RunE:  createDeploymentExecution,
}

func init() {
	CreateDeployment.Flags().StringVarP(&storageRepoURI, "storage-repo-uri", "s", "", "storage git repo uri")
	CreateDeployment.MarkFlagRequired("storage-repo-uri")

	CreateDeployment.Flags().StringVarP(&codeRepoURI, "code-repo-uri", "c", "", "code git repo uri")
	CreateDeployment.MarkFlagRequired("code-repo-uri")

	CreateDeployment.Flags().StringVarP(&codeRepoPath, "code-repo-path", "p", "", "code git repo path")
	CreateDeployment.Flags().StringVarP(&terraformVersion, "terraform-version", "t", "", "terraform version")
	CreateDeployment.Flags().StringVarP(&flavour, "flavour", "f", "", "flavour")
}

func createDeploymentExecution(command *cobra.Command, args []string) error {
	deployName := args[0]
	m := manager.GetManager()

	// Configurable default values for flags
	if flavour == "" {
		flavour = viper.GetString("DefaultFlavour")
	}
	if terraformVersion == "" {
		terraformVersion = viper.GetString("DefaultTerraformVersion")
	}

	err := m.Create(deployName, storageRepoURI, codeRepoURI, codeRepoPath, terraformVersion, flavour)
	if err != nil {
		return err
	}

	err = common.SetCurrentDeployment(deployName)
	if err != nil {
		return err
	}

	fmt.Println("Created")
	return nil
}
