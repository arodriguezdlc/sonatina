package operation

import (
	"github.com/spf13/cobra"
)

// Delete declares `sonatina delete` command
var Delete = &cobra.Command{
	Use:   "delete",
	Short: "d",
	Long:  "delete resources",
	//Run:   util.Help,
}

func init() {
	//createCmd.AddCommand(api.CreateCmd)
}
