package deploymentcmd

import (
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/arodriguezdlc/sonatina/manager"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ListDeployment declares `sonatina list deployments` command
var ListDeployment = &cobra.Command{
	Use:   "deployments",
	Short: "list deployments managed by sonatina",
	Long:  `TO DO`,
	Run:   listDeploymentExecution,
}

//To define flags
func init() {
}

func listDeploymentExecution(cmd *cobra.Command, args []string) {
	m := manager.GetManager()
	list, err := m.List()

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
