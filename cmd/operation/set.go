package operation

import (
	"github.com/arodriguezdlc/sonatina/cmd/flavour"
	"github.com/spf13/cobra"
)

// Set declares `sonatina set` command
var Set = &cobra.Command{
	Use:   "set",
	Short: "",
	Long:  "",
}

func init() {
	Set.AddCommand(flavour.SetFlavour)
}
