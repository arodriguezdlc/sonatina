package deployment

import (
	"path/filepath"

	"github.com/arodriguezdlc/sonatina/gitw"
	"github.com/arodriguezdlc/sonatina/utils"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// TODO: improve docs

const varsBranch string = "variables"

// Vars manage variable and metadata
type Vars struct {
	fs   afero.Fs
	path string
	gitw *gitw.Command

	RepoURL  string
	Metadata *Metadata
}

// getVars creates and initializes a new Vars object
func (d *DeploymentImpl) getVars(repoURL string) error {
	vars, err := newVars(d.fs, d.path, repoURL)
	if err != nil {
		return err
	}

	err = vars.Metadata.load()
	if err != nil {
		return err
	}

	d.Vars = vars
	return nil
}

// cloneVars creates and initializes a new Vars object
func (d *DeploymentImpl) cloneVars(repoURL string) error {
	vars, err := newVars(d.fs, d.path, repoURL)
	if err != nil {
		return err
	}

	err = vars.gitw.CloneBranch(repoURL, varsBranch)
	if err != nil {
		return err
	}

	err = vars.Metadata.load()
	if err != nil {
		return err
	}

	d.Vars = vars
	return nil
}

func (d *DeploymentImpl) createVars(storageRepoURL string, terraformVersion string, codeRepoURL string, codeRepoPath string, flavour string) error {
	vars, err := newVars(d.fs, d.path, storageRepoURL)
	if err != nil {
		return err
	}

	vars.Metadata.TerraformVersion = terraformVersion
	vars.Metadata.Repo = codeRepoURL
	vars.Metadata.RepoPath = codeRepoPath
	vars.Metadata.Flavour = flavour
	err = vars.Metadata.save()
	if err != nil {
		return err
	}

	err = vars.gitw.Init()
	if err != nil {
		return err
	}

	err = vars.gitw.RemoteAdd("origin", storageRepoURL)
	if err != nil {
		return err
	}

	for _, subdir := range []string{"global", "user"} {
		path := filepath.Join(vars.path, subdir)
		err = utils.NewDirectoryWithKeep(vars.fs, path)
		if err != nil {
			return err
		}
	}

	err = vars.gitw.AddGlob(".")
	if err != nil {
		return err
	}

	err = vars.gitw.Commit("Initial commit")
	if err != nil {
		return err
	}

	// XXX: checkout executed after first commit to avoid
	// reference errors. Must be reviewed
	err = vars.gitw.CheckoutNewBranch(varsBranch)
	if err != nil {
		return err
	}

	d.Vars = vars
	return nil
}

func newVars(fs afero.Fs, deploymentPath string, repoURL string) (*Vars, error) {
	path := filepath.Join(deploymentPath, varsBranch)

	err := fs.MkdirAll(path, 0755)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create directory")
	}

	varsGit, err := gitw.NewCommand(fs, path)
	if err != nil {
		return nil, err
	}

	vars := &Vars{
		fs:   fs,
		path: path,
		gitw: varsGit,

		RepoURL:  repoURL,
		Metadata: newMetadata(fs, path),
	}

	return vars, nil
}
