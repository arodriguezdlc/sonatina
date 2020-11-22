package operation

import (
	"github.com/arodriguezdlc/sonatina/cmd/deploymentcmd"
	"github.com/arodriguezdlc/sonatina/cmd/plugin"
	"github.com/arodriguezdlc/sonatina/cmd/usercomponent"
	"github.com/spf13/cobra"
)

// List declares `sonatina list` command
var List = &cobra.Command{
	Use:   "list",
	Short: "List a set of resources",
}

func init() {
	List.AddCommand(deploymentcmd.ListDeployment)
	List.AddCommand(usercomponent.ListUsercomponents)
	List.AddCommand(plugin.ListPlugins)
}
