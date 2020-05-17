package operation

import (
	"github.com/arodriguezdlc/sonatina/cmd/deploymentcmd"
	"github.com/spf13/cobra"
)

// Clone declares `sonatina clone` command
var Clone = &cobra.Command{
	Use:   "clone",
	Short: "",
	Long:  "",
}

func init() {
	Clone.AddCommand(deploymentcmd.CloneDeployment)
}
