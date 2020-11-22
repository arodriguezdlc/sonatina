package operation

import (
	"github.com/arodriguezdlc/sonatina/cmd/deploymentcmd"
	"github.com/arodriguezdlc/sonatina/cmd/plugin"
	"github.com/arodriguezdlc/sonatina/cmd/usercomponent"
	"github.com/spf13/cobra"
)

// Delete declares `sonatina delete` command
var Delete = &cobra.Command{
	Use:   "delete",
	Short: "Remove a resource",
}

func init() {
	Delete.AddCommand(deploymentcmd.DeleteDeployment)
	Delete.AddCommand(usercomponent.DeleteUsercomponent)
	Delete.AddCommand(plugin.DeletePlugin)
}
