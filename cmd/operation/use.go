package operation

import (
	"github.com/arodriguezdlc/sonatina/cmd/deploymentcmd"
	"github.com/spf13/cobra"
)

// Use declares `sonatina use` command
var Use = &cobra.Command{
	Use:   "use",
	Short: "Configure a default deployment to do operations",
}

func init() {
	Use.AddCommand(deploymentcmd.UseDeployment)
}
