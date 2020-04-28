package deployment

import (
	"github.com/arodriguezdlc/sonatina/gitw"
	"github.com/arodriguezdlc/sonatina/utils"
	"github.com/spf13/afero"
)

// TODO: improve docs

// Vars manage variable and metadata
type Vars struct {
	fs       afero.Fs
	path     string
	repoURL  string
	gitw     gitw.Command
	metadata metadata
}

// NewVars creates and initializes a new Vars object
func NewVars(fs afero.Fs, path string, repoURL string) (*Vars, error) {
	var err error

	vars := &Vars{
		fs:      fs,
		path:    path,
		repoURL: repoURL,
		gitw:    gitw.Command{},
	}

	vars.gitw, err = gitw.NewCommand(fs, path)
	if err != nil {
		return nil, err
	}

	ok, err := vars.isInitialized()
	if err != nil {
		return nil, err
	}
	if !ok {
		operation, err := vars.gitw.CloneOrInit(repoURL, "variables")
		if err != nil {
			return nil, err
		}
		if operation == "init" {
			if err = vars.initialize(); err != nil {
				return nil, err
			}
		}
	}

	return vars, nil
}

// Save method stores terraform vars information on git repository
func (v *Vars) Save() {
	// TODO
}

func (v *Vars) initialize() error {
	// TODO: add metadata
	err := utils.NewFileIfNotExist(v.fs, v.path+"/metadata.yml")
	if err != nil {
		return err
	}

	err = v.fs.MkdirAll(v.path+"/global", 0755)
	if err != nil {
		return err
	}

	err = utils.NewFileIfNotExist(v.fs, v.path+"/global/.keep")
	if err != nil {
		return err
	}

	err = v.gitw.AddGlob(".")
	if err != nil {
		return err
	}

	err = v.gitw.Commit("Initial commit")
	if err != nil {
		return err
	}

	return v.gitw.CheckoutNewBranch("variables")
}

func (v *Vars) isInitialized() (bool, error) {
	return afero.Exists(v.fs, v.path)
}
