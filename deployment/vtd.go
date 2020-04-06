package deployment

import "github.com/spf13/afero"

// VTD represents a Variable Tree Definition
type VTD struct {
	fs   afero.Fs
	path string
}

func (vtd *VTD) ListConfigGlobal() ([]string, error) {
	return afero.Glob(vtd.fs, vtd.path+"/config/global/*.tfvars")
}

func (vtd *VTD) ListConfigUser(user string) ([]string, error) {
	return afero.Glob(vtd.fs, vtd.path+"/config/user/"+user+"/*.tfvars")
}

func (vtd *VTD) ListFlavourGlobal(flavour string) ([]string, error) {
	return afero.Glob(vtd.fs, vtd.path+"/flavour/global/"+flavour+".tfvars")
}

func (vtd *VTD) ListFlavourUser(flavour string, user string) ([]string, error) {
	return afero.Glob(vtd.fs, vtd.path+"/flavour/user/"+user+"/"+flavour+".tfvars")
}

func (vtd *VTD) ListStaticGlobal() ([]string, error) {
	return afero.Glob(vtd.fs, vtd.path+"/static/global/*.tfvars")
}

func (vtd *VTD) ListStaticUser(user string) ([]string, error) {
	return afero.Glob(vtd.fs, vtd.path+"/static/user/"+user+"/*.tfvars")
}
