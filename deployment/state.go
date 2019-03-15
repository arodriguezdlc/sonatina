package deployment

// State manages terraform state
type State interface {
	Path() string
	Save()
	Load()
}
