package deploymentcmd

import (
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/arodriguezdlc/sonatina/manager"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// AddDeployment declares `sonatina add deployment` command
var AddDeployment = &cobra.Command{
	Use:   "deployment",
	Short: "deployment",
	Long:  "deployment",
	Run:   addDeploymentExecution,
}

//To define flags
func init() {
}

func addDeploymentExecution(cmd *cobra.Command, args []string) {
	m := manager.GetManager()
	list, err := manager.Add()

	if err != nil {
		log.Fatalln(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Deployments"})
	for _, element := range list {
		var aux []string
		aux = append(aux, element)
		table.Append(aux)
	}
	table.Render()
}
