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
	fs         afero.Fs
	path       string
	gitw       *gitw.Command
	deployment *DeploymentImpl

	RepoURL  string
	Metadata *Metadata
}

// CreateUsercomponent adds a new user component to metadata and
// initializes the directory tree
func (v *Vars) CreateUsercomponent(user string) error {
	err := v.fs.MkdirAll(v.UsercomponentPath(user), 0755)
	if err != nil {
		return errors.Wrap(err, "couldn't create directory")
	}

	// TODO: Create user component on metadata
	err = v.Metadata.CreateUsercomponent(user)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUsercomponent deletes an user component from metadata
// and performs a directory cleanup
func (v *Vars) DeleteUsercomponent(user string) error {
	err := v.Metadata.DeleteUsercomponent(user)
	if err != nil {
		return err
	}

	err = v.fs.RemoveAll(v.UsercomponentPath(user))
	if err != nil {
		return errors.Wrap(err, "couldn't remove dir recursively")
	}

	return nil
}

// UsercomponentPath returns de variable directory path for a
// specified user
func (v *Vars) UsercomponentPath(user string) string {
	return filepath.Join(v.path, "user", user)
}

// GenerateGlobal generates vars files to be used on
// terraform operations. Returns a list of vars files that
// must be applied in order.
func (v *Vars) GenerateGlobal() ([]string, error) {
	varFiles, err := v.copyVTDGlobal(v.deployment.Base.vtd, "base", v.Metadata.Flavour)
	if err != nil {
		return varFiles, err
	}

	for i, plugin := range v.deployment.Plugins {
		pluginFiles, err := v.copyVTDGlobal(plugin.vtd, "plugin_"+v.Metadata.Plugins[i].Name, v.Metadata.Flavour)
		if err != nil {
			return varFiles, err
		}
		varFiles = append(varFiles, pluginFiles...)
	}

	return varFiles, nil
}

// GenerateUser generates vars files to be used on
// terraform operations. Returns a list of vars files that
// must be applied in order.
func (v *Vars) GenerateUser(user string) ([]string, error) {
	varFiles, err := v.copyVTDUser(user, v.deployment.Base.vtd, "base", v.Metadata.UserComponents[user].Flavour)
	if err != nil {
		return varFiles, err
	}

	userPluginList, err := v.Metadata.ListUserPlugins(user)
	if err != nil {
		return varFiles, err
	}

	for i, plugin := range v.deployment.Plugins {
		// Only adds plugin if user has the plugin added.
		_, ok := utils.FindString(userPluginList, v.Metadata.Plugins[i].Name)
		if ok {
			pluginFiles, err := v.copyVTDUser(user, plugin.vtd, "plugin_"+v.Metadata.Plugins[i].Name, v.Metadata.UserComponents[user].Flavour)
			if err != nil {
				return varFiles, err
			}
			varFiles = append(varFiles, pluginFiles...)
		}
	}

	return varFiles, nil
}

func (v *Vars) copyVTDGlobal(vtd *VTD, prefix string, flavour string) ([]string, error) {
	varFiles := []string{}

	staticFile, err := v.copyStaticGlobal(vtd, prefix)
	if err != nil {
		return varFiles, err
	}

	flavourFile, err := v.copyFlavourGlobal(vtd, prefix, flavour)
	if err != nil {
		return varFiles, err
	}

	configFile, err := v.copyConfigGlobal(vtd, prefix)
	if err != nil {
		return varFiles, err
	}

	varFiles = []string{staticFile, flavourFile, configFile}
	return varFiles, nil
}

func (v *Vars) copyVTDUser(user string, vtd *VTD, prefix string, flavour string) ([]string, error) {
	varFiles := []string{}

	staticFile, err := v.copyStaticUser(user, vtd, prefix)
	if err != nil {
		return varFiles, err
	}

	flavourFile, err := v.copyFlavourUser(user, vtd, prefix, flavour)
	if err != nil {
		return varFiles, err
	}

	configFile, err := v.copyConfigUser(user, vtd, prefix)
	if err != nil {
		return varFiles, err
	}

	varFiles = []string{staticFile, flavourFile, configFile}
	return varFiles, nil
}

func (v *Vars) copyConfigGlobal(vtd *VTD, prefix string) (string, error) {
	src := vtd.config.globalFile()
	dst := filepath.Join(v.path, "global", prefix+"_config.tfvars")

	ok, err := afero.Exists(v.fs, dst)
	if err != nil {
		return dst, errors.Wrap(err, "couldn't determine if file exists")
	}
	if !ok {
		err := utils.FileCopy(v.fs, src, dst)
		if err != nil {
			return dst, err
		}
	}

	return dst, nil
}

func (v *Vars) copyConfigUser(user string, vtd *VTD, prefix string) (string, error) {
	src := vtd.config.userFile()
	dst := filepath.Join(v.path, "user", user, prefix+"_config.tfvars")

	ok, err := afero.Exists(v.fs, dst)
	if err != nil {
		return dst, errors.Wrap(err, "couldn't determine if file exists")
	}
	if !ok {
		err := utils.FileCopy(v.fs, src, dst)
		if err != nil {
			return dst, err
		}
	}

	return dst, nil
}

func (v *Vars) copyFlavourGlobal(vtd *VTD, prefix string, flavour string) (string, error) {
	src := vtd.flavour.globalFile(flavour)
	dst := filepath.Join(v.path, "global", prefix+"_flavour_"+flavour+".tfvars")

	err := utils.FileCopy(v.fs, src, dst)
	if err != nil {
		return dst, err
	}

	return dst, nil
}

func (v *Vars) copyFlavourUser(user string, vtd *VTD, prefix string, flavour string) (string, error) {
	src := vtd.flavour.userFile(flavour)
	dst := filepath.Join(v.path, "user", user, prefix+"_flavour_"+flavour+".tfvars")

	err := utils.FileCopy(v.fs, src, dst)
	if err != nil {
		return dst, err
	}

	return dst, nil
}

func (v *Vars) copyStaticGlobal(vtd *VTD, prefix string) (string, error) {
	src := vtd.static.globalFile()
	dst := filepath.Join(v.path, "global", prefix+"_static.tfvars")

	err := utils.FileCopy(v.fs, src, dst)
	if err != nil {
		return dst, err
	}

	return dst, nil
}

func (v *Vars) copyStaticUser(user string, vtd *VTD, prefix string) (string, error) {
	src := vtd.static.globalFile()
	dst := filepath.Join(v.path, "user", user, prefix+"_static.tfvars")

	err := utils.FileCopy(v.fs, src, dst)
	if err != nil {
		return dst, err
	}

	return dst, nil
}

// getVars creates and initializes a new Vars object
func (d *DeploymentImpl) getVars(repoURL string) error {
	vars, err := d.newVars(repoURL)
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
	vars, err := d.newVars(repoURL)
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
	vars, err := d.newVars(storageRepoURL)
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

func (d *DeploymentImpl) newVars(repoURL string) (*Vars, error) {
	path := filepath.Join(d.path, varsBranch)

	err := d.fs.MkdirAll(path, 0755)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create directory")
	}

	varsGit, err := gitw.NewCommand(d.fs, path)
	if err != nil {
		return nil, err
	}

	vars := &Vars{
		fs:         d.fs,
		path:       path,
		gitw:       varsGit,
		deployment: d,

		RepoURL:  repoURL,
		Metadata: newMetadata(d.fs, path),
	}

	return vars, nil
}
