package operation

import (
	"github.com/arodriguezdlc/sonatina/cmd/deploymentcmd"
	"github.com/spf13/cobra"
)

// Add declares `sonatina add` command
var Add = &cobra.Command{
	Use:   "add",
	Short: "",
	Long:  "",
	//Run:   util.Help,
}

func init() {
	Add.AddCommand(deploymentcmd.AddDeployment)
}
