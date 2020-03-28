package manager

import (
	"reflect"
	"testing"

	"github.com/spf13/afero"
)

func TestListWithCorrectFile(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := testWriteCorrectDeployFile(fs)
	if err != nil {
		t.Error(err)
	}

	m, err := newManagerJSON(fs, "/", "deployments.json")
	if err != nil {
		t.Fatal(err)
	}

	result, err := m.List()
	if err != nil {
		t.Fatal(err)
	}

	expectedArray := [2]string{"deploy1", "deploy2"}
	expected := expectedArray[:]
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Incorrect array. Expected: %s, Obtained: %s", expected, result)
	}
}

func TestListWithEmptyFile(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := testWriteEmptyDeployFile(fs)
	if err != nil {
		t.Error(err)
	}

	m, err := newManagerJSON(fs, "/", "deployments.json")
	if err != nil {
		t.Fatal(err)
	}

	result, err := m.List()
	if err != nil {
		t.Fatal(err)
	}

	expectedArray := [0]string{}
	expected := expectedArray[:]
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Incorrect array. Expected: %s, Obtained: %s", expected, result)
	}
}

// TODO: use mocks
// func TestGetDeploy(t *testing.T) {
// 	fs := afero.NewMemMapFs()
// 	deploymentFile := testWriteCorrectDeployFile(fs)

// 	m, err := newManagerJSON(fs, "/", deploymentFile)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	result, err := m.Get("deploy1")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	expectedVars, err := deployment.NewVars(fs, "/variables", "git@test.com/deploy1storage")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	expectedState, err := deployment.NewState(fs, "/state", "git@test.com/deploy1storage")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	expectedDeploy := &deployment.DeploymentImpl{
// 		Name:  "deploy1",
// 		Vars:  expectedVars,
// 		State: expectedState,
// 	}

// 	if !reflect.DeepEqual(result, *expectedDeploy) {
// 		t.Errorf("Incorrect deploy. Expected: %v, Obtained: %v", *expectedDeploy, result)
// 	}
// }

func testWriteCorrectDeployFile(fs afero.Fs) error {
	return afero.WriteFile(fs, "/deployments.json", []byte(testCorrectDeployContent()), 0644)
}

func testWriteEmptyDeployFile(fs afero.Fs) error {
	return afero.WriteFile(fs, "/deployments.json", []byte("{}"), 0644)
}

func testCorrectDeployContent() string {
	return `{
    "deploy1": {
      "storage_repo_uri": "git@test.com/deploy1storage",
      "code_repo_uri": "git@test.com/deploy1code"
    },
    "deploy2": {
      "storage_repo_uri": "git@test.com/deploy2storage",
      "code_repo_uri": "git@test.com/deploy2code"
    }
	}`
}
