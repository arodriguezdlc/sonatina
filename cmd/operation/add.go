package operation

import (
	"github.com/spf13/cobra"
)

// Add declares `sonatina add` command
var Add = &cobra.Command{
	Use:   "add",
	Short: "a",
	Long:  "add",
	//Run:   util.Help,
}

func init() {
	//createCmd.AddCommand(api.CreateCmd)
}
