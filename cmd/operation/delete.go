package operation

import (
	"github.com/arodriguezdlc/sonatina/cmd/deploymentcmd"
	"github.com/spf13/cobra"
)

// Delete declares `sonatina delete` command
var Delete = &cobra.Command{
	Use:   "delete",
	Short: "delete",
	Long:  "delete",
	//Run:   util.Help,
}

func init() {
	Delete.AddCommand(deploymentcmd.DeleteDeployment)
}
