package operation

import (
	"github.com/arodriguezdlc/sonatina/cmd/flavour"
	"github.com/spf13/cobra"
)

// Get declares `sonatina get` command
var Get = &cobra.Command{
	Use:   "get",
	Short: "Obtain information about a resource",
}

func init() {
	Get.AddCommand(flavour.GetFlavour)
}
