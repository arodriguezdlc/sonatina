package workflow

// TODO: improve godoc

// Workflow interface
type Workflow interface {
	// TODO: define a workflow interface to be able to implement operations
	// over all workflows

	Run() error
}
