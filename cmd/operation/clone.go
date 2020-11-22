package operation

import (
	"github.com/arodriguezdlc/sonatina/cmd/deploymentcmd"
	"github.com/spf13/cobra"
)

// Clone declares `sonatina clone` command
var Clone = &cobra.Command{
	Use:   "clone",
	Short: "Get resources from git",
}

func init() {
	Clone.AddCommand(deploymentcmd.CloneDeployment)
}
