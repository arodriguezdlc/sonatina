package operation

import (
	"github.com/arodriguezdlc/sonatina/cmd/deploymentcmd"
	"github.com/arodriguezdlc/sonatina/cmd/usercomponent"
	"github.com/spf13/cobra"
)

// List declares `sonatina list` command
var List = &cobra.Command{
	Use:   "list",
	Short: "list",
	Long:  "list",
	//Run:   util.Help,
}

func init() {
	List.AddCommand(deploymentcmd.ListDeployment)
	List.AddCommand(usercomponent.ListUsercomponents)
}
