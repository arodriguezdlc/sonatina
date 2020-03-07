package deployment

import (
	"github.com/arodriguezdlc/sonatina/gitw"
	"github.com/arodriguezdlc/sonatina/utils"
	"github.com/spf13/afero"
)

// VarsGit implements Vars interface
type VarsGit struct {
	fs      afero.Fs
	path    string
	repoURL string
	gitw    gitw.Command
}

// NewVarsGit creates and initializes a new VarsGit object
func NewVarsGit(fs afero.Fs, path string, repoURL string) (Vars, error) {
	var err error

	vars := VarsGit{
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

// Path returns the path where vars files are saved
func (v VarsGit) Path() string {
	return v.path
}

// Save method stores terraform vars information on git repository
func (v VarsGit) Save() {
	// TODO
}

func (v VarsGit) initialize() error {
	// TODO: add metadata
	err := utils.NewFileIfNotExist(v.path+"/metadata.yml", v.fs)
	if err != nil {
		return err
	}

	err = v.fs.MkdirAll(v.path+"/global", 0700)
	if err != nil {
		return err
	}

	err = utils.NewFileIfNotExist(v.path+"/global/.keep", v.fs)
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

func (v VarsGit) isInitialized() (bool, error) {
	return afero.Exists(v.fs, v.path)
}
