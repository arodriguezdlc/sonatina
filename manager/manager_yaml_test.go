package manager

import (
	"reflect"
	"testing"

	"github.com/arodriguezdlc/sonatina/deployment"

	"github.com/spf13/afero"
)

func TestListWithCorrectFile(t *testing.T) {
	var fs afero.Fs
	var deploymentFile string
	var m ManagerYaml
	var err error
	var result deployment.Deployment

	fs = afero.NewMemMapFs()
	deploymentFile = writeCorrectDeployFile(fs)

	if m, err = NewManagerYaml(fs, "/", deploymentFile); err != nil {
		t.Fatal(err)
	}

	if result, err = m.List(); err != nil {
		t.Fatal(err)
	}

	expectedArray := [2]string{"deploy1", "deploy2"}
	expected := expectedArray[:]
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Incorrect array. Expected: %s, Obtained: %s", expected, result)
	}
}

func TestGetDeploy(t *testing.T) {
	var fs afero.Fs
	var deploymentFile string
	var m ManagerYaml
	var err error
	var result deployment.Deployment
	var expectedVars deployment.Vars
	var expectedState deployment.State
	var expectedDeploy deployment.Deployment

	fs = afero.NewMemMapFs()
	deploymentFile = writeCorrectDeployFile(fs)

	if m, err = NewManagerYaml(fs, "/", deploymentFile); err != nil {
		t.Fatal(err)
	}

	if result, err = m.Get("deploy1"); err != nil {
		t.Fatal(err)
	}

	if expectedVars, err = deployment.NewVarsGit(fs, "/variables", "git@test.com/deploy1storage"); err != nil {
		t.Fatal(err)
	}

	if expectedState, err = deployment.NewStateGit(fs, "/state", "git@test.com/deploy1storage"); err != nil {
		t.Fatal(err)
	}

	expectedDeploy = deployment.DeploymentImpl{
		Name:  "deploy1",
		Vars:  expectedVars,
		State: expectedState,
	}

	if !reflect.DeepEqual(result, expectedDeploy) {
		t.Errorf("Incorrect deploy. Expected: %s, Obtained: %s", expectedDeploy, result)
	}
}

func writeCorrectDeployFile(fs afero.Fs) string {
	deploymentFile := "/deployments.yml"
	content := `---
deploy1:
  storage_repo_uri: "git@test.com/deploy1storage"
  code_repo_uri: "git@test.com/deploy1code"
deploy2:
  storage_repo_uri: "git@test.com/deploy2storage"
  code_repo_uri: "git@test.com/deploy2code"
`
	afero.WriteFile(fs, deploymentFile, []byte(content), 0644)
	return deploymentFile
}
