package operation

import (
	"github.com/arodriguezdlc/sonatina/cmd/deploymentcmd"
	"github.com/arodriguezdlc/sonatina/cmd/usercomponent"
	"github.com/spf13/cobra"
)

// Create declares `sonatina create` command
var Create = &cobra.Command{
	Use:   "create",
	Short: "",
	Long:  "",
}

func init() {
	Create.AddCommand(deploymentcmd.CreateDeployment)
	Create.AddCommand(usercomponent.CreateUsercomponent)
}
