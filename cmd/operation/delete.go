package operation

import (
	"github.com/arodriguezdlc/sonatina/cmd/deploymentcmd"
	"github.com/arodriguezdlc/sonatina/cmd/usercomponent"
	"github.com/spf13/cobra"
)

// Delete declares `sonatina delete` command
var Delete = &cobra.Command{
	Use:   "delete",
	Short: "delete",
	Long:  "delete",
}

func init() {
	Delete.AddCommand(deploymentcmd.DeleteDeployment)
	Delete.AddCommand(usercomponent.DeleteUsercomponent)
}
