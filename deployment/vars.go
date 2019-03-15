package deployment

// Vars manage variables
type Vars interface {
	Path() string
	Save()
	Load()
}
