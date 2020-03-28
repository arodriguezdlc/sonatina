package workflow

import "github.com/arodriguezdlc/sonatina/deployment"

type ProcessCTD struct {
	deployment deployment.Deployment
}

func (p *ProcessCTD) Run() {
	// generate file list to copy to workdir
	// consult files from all CTDs using its methods

	// Copy all files to workdir applying override rules.
	// Modules are copied too, and also are overrided.

}

func (p *ProcessCTD) generateMainFileList() {

}

func (p *ProcessCTD) generateModuleList() {

}
