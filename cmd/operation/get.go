package operation

import (
	"github.com/spf13/cobra"
)

// Get declares `sonatina add` command
var Get = &cobra.Command{
	Use:   "get",
	Short: "get",
	Long:  "get",
	//Run:   util.Help,
}

func init() {
	//Get.AddCommand(deploymentcmd.GetDeployment)
}
