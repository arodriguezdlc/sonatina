package workflow

import "github.com/arodriguezdlc/sonatina/deployment"

// should implements Workflow interface
type Init struct {
	deployment deployment.Deployment
}

func (i *Init) Run() {

}
