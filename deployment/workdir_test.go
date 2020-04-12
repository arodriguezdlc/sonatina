package deployment

import (
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/arodriguezdlc/sonatina/utils"
	"github.com/spf13/afero"
)

func TestGenerateGlobal(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := testWordirCreateDeploymentDirectories(fs)
	if err != nil {
		t.Fatal(err)
	}

	workdir, err := testNewWorkdir(fs)
	if err != nil {
		t.Fatal(err)
	}

	err = workdir.GenerateGlobal()
	if err != nil {
		t.Fatal(err)
	}

	expectedFileList := testWorkdirGlobalExpectedFiles()
	obtainedFileList, err := utils.FileListRecursivelyWithoutDirs(fs, filepath.Join("deployment", "workdir"))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedFileList, obtainedFileList) {
		t.Errorf("Incorrect file list.\n\n Expected: %v\n\n Obtained: %v\n", expectedFileList, obtainedFileList)
	}
}

func TestGenerateGlobalWithOverride(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := testWordirCreateDeploymentDirectoriesWithOverride(fs)
	if err != nil {
		t.Fatal(err)
	}

	workdir, err := testNewWorkdir(fs)
	if err != nil {
		t.Fatal(err)
	}

	err = workdir.GenerateGlobal()
	if err != nil {
		t.Fatal(err)
	}

	expectedFileList := testWorkdirGlobalExpectedFilesWithOverride()
	obtainedFileList, err := utils.FileListRecursivelyWithoutDirs(fs, filepath.Join("deployment", "workdir"))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedFileList, obtainedFileList) {
		t.Errorf("Incorrect file list.\n\n Expected: %v\n\n Obtained: %v\n", expectedFileList, obtainedFileList)
	}

}

func TestGenerateUser(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := testWordirCreateDeploymentDirectories(fs)
	if err != nil {
		t.Fatal(err)
	}

	workdir, err := testNewWorkdir(fs)
	if err != nil {
		t.Fatal(err)
	}

	err = workdir.GenerateUser("user1")
	if err != nil {
		t.Fatal(err)
	}

	expectedFileList := testWorkdirUserExpectedFiles("user1")
	obtainedFileList, err := utils.FileListRecursivelyWithoutDirs(fs, filepath.Join("deployment", "workdir"))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedFileList, obtainedFileList) {
		t.Errorf("Incorrect file list.\n\n Expected: %v\n\n Obtained: %v\n", expectedFileList, obtainedFileList)
	}

	err = workdir.GenerateUser("user2")
	if err != nil {
		t.Fatal(err)
	}

	expectedFileListWithDuplicates := append(expectedFileList, testWorkdirUserExpectedFiles("user2")...)
	expectedFileList = utils.RemoveDuplicatedStrings(expectedFileListWithDuplicates)
	sort.Strings(expectedFileList)
	obtainedFileList, err = utils.FileListRecursivelyWithoutDirs(fs, filepath.Join("deployment", "workdir"))
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(obtainedFileList)

	if !reflect.DeepEqual(expectedFileList, obtainedFileList) {
		t.Errorf("Incorrect file list.\n\n Expected: %v\n\n Obtained: %v\n", expectedFileList, obtainedFileList)
	}
}

func TestGenerateUserWithOverride(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := testWordirCreateDeploymentDirectoriesWithOverride(fs)
	if err != nil {
		t.Fatal(err)
	}

	workdir, err := testNewWorkdir(fs)
	if err != nil {
		t.Fatal(err)
	}

	err = workdir.GenerateUser("user1")
	if err != nil {
		t.Fatal(err)
	}

	expectedFileList := testWorkdirUserExpectedFilesWithOverride("user1")
	obtainedFileList, err := utils.FileListRecursivelyWithoutDirs(fs, filepath.Join("deployment", "workdir"))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedFileList, obtainedFileList) {
		t.Errorf("Incorrect file list.\n\n Expected: %v\n\n Obtained: %v\n", expectedFileList, obtainedFileList)
	}
}

func testNewWorkdir(fs afero.Fs) (*Workdir, error) {

	base, err := NewCTD(fs, filepath.Join("deployment", "base"), "", "")
	if err != nil {
		return &Workdir{}, err
	}

	plugin1, err := NewCTD(fs, filepath.Join("deployment", "plugins", "plugin1"), "", "")
	if err != nil {
		return &Workdir{}, err
	}

	plugin2, err := NewCTD(fs, filepath.Join("deployment", "plugins", "plugin2"), "", "")
	if err != nil {
		return &Workdir{}, err
	}

	workdirCTD, err := NewCTD(fs, filepath.Join("deployment", "workdir"), "", "")
	if err != nil {
		return &Workdir{}, err
	}

	deploy := &DeploymentImpl{
		Name: "deployment",
		fs:   fs,
		path: "/deployment",

		Base:    base,
		Plugins: []*CTD{plugin1, plugin2},
	}

	workdir := &Workdir{
		fs:         fs,
		path:       "/deployment/workdir",
		deployment: deploy,
		CTD:        workdirCTD,
	}

	return workdir, nil
}

func testWordirCreateDeploymentDirectories(fs afero.Fs) error {
	basepaths := []string{
		filepath.Join("deployment", "base"),
		filepath.Join("deployment", "plugins", "plugin1"),
		filepath.Join("deployment", "plugins", "plugin2"),
	}

	for _, path := range basepaths {
		err := testWorktreeCreateCTDFiles(fs, path, filepath.Base(path)+"_")
		if err != nil {
			return err
		}
	}

	return nil
}

func testWordirCreateDeploymentDirectoriesWithOverride(fs afero.Fs) error {
	basepaths := []string{
		filepath.Join("deployment", "base"),
		filepath.Join("deployment", "plugins", "plugin1"),
		filepath.Join("deployment", "plugins", "plugin2"),
	}

	for _, path := range basepaths {
		err := testWorktreeCreateCTDFiles(fs, path, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func testWorkdirGlobalExpectedFiles() []string {
	return []string{
		filepath.Join("deployment", "workdir", "main", "global", "base_file1.tf"),
		filepath.Join("deployment", "workdir", "main", "global", "base_file2.tf"),
		filepath.Join("deployment", "workdir", "main", "global", "plugin1_file1.tf"),
		filepath.Join("deployment", "workdir", "main", "global", "plugin1_file2.tf"),
		filepath.Join("deployment", "workdir", "main", "global", "plugin2_file1.tf"),
		filepath.Join("deployment", "workdir", "main", "global", "plugin2_file2.tf"),
		filepath.Join("deployment", "workdir", "modules", "base_module1", "file1.tf"),
		filepath.Join("deployment", "workdir", "modules", "base_module1", "file2.tf"),
		filepath.Join("deployment", "workdir", "modules", "base_module2", "file1.tf"),
		filepath.Join("deployment", "workdir", "modules", "base_module2", "file2.tf"),
		filepath.Join("deployment", "workdir", "modules", "plugin1_module1", "file1.tf"),
		filepath.Join("deployment", "workdir", "modules", "plugin1_module1", "file2.tf"),
		filepath.Join("deployment", "workdir", "modules", "plugin1_module2", "file1.tf"),
		filepath.Join("deployment", "workdir", "modules", "plugin1_module2", "file2.tf"),
		filepath.Join("deployment", "workdir", "modules", "plugin2_module1", "file1.tf"),
		filepath.Join("deployment", "workdir", "modules", "plugin2_module1", "file2.tf"),
		filepath.Join("deployment", "workdir", "modules", "plugin2_module2", "file1.tf"),
		filepath.Join("deployment", "workdir", "modules", "plugin2_module2", "file2.tf"),
	}
}

func testWorkdirGlobalExpectedFilesWithOverride() []string {
	return []string{
		filepath.Join("deployment", "workdir", "main", "global", "file1.tf"),
		filepath.Join("deployment", "workdir", "main", "global", "file2.tf"),
		filepath.Join("deployment", "workdir", "modules", "module1", "file1.tf"),
		filepath.Join("deployment", "workdir", "modules", "module1", "file2.tf"),
		filepath.Join("deployment", "workdir", "modules", "module2", "file1.tf"),
		filepath.Join("deployment", "workdir", "modules", "module2", "file2.tf"),
	}
}

func testWorkdirUserExpectedFiles(user string) []string {
	return []string{
		filepath.Join("deployment", "workdir", "main", "user", user, "base_file1.tf"),
		filepath.Join("deployment", "workdir", "main", "user", user, "base_file2.tf"),
		filepath.Join("deployment", "workdir", "main", "user", user, "plugin1_file1.tf"),
		filepath.Join("deployment", "workdir", "main", "user", user, "plugin1_file2.tf"),
		filepath.Join("deployment", "workdir", "main", "user", user, "plugin2_file1.tf"),
		filepath.Join("deployment", "workdir", "main", "user", user, "plugin2_file2.tf"),
		filepath.Join("deployment", "workdir", "modules", "base_module1", "file1.tf"),
		filepath.Join("deployment", "workdir", "modules", "base_module1", "file2.tf"),
		filepath.Join("deployment", "workdir", "modules", "base_module2", "file1.tf"),
		filepath.Join("deployment", "workdir", "modules", "base_module2", "file2.tf"),
		filepath.Join("deployment", "workdir", "modules", "plugin1_module1", "file1.tf"),
		filepath.Join("deployment", "workdir", "modules", "plugin1_module1", "file2.tf"),
		filepath.Join("deployment", "workdir", "modules", "plugin1_module2", "file1.tf"),
		filepath.Join("deployment", "workdir", "modules", "plugin1_module2", "file2.tf"),
		filepath.Join("deployment", "workdir", "modules", "plugin2_module1", "file1.tf"),
		filepath.Join("deployment", "workdir", "modules", "plugin2_module1", "file2.tf"),
		filepath.Join("deployment", "workdir", "modules", "plugin2_module2", "file1.tf"),
		filepath.Join("deployment", "workdir", "modules", "plugin2_module2", "file2.tf"),
	}
}

func testWorkdirUserExpectedFilesWithOverride(user string) []string {
	return []string{
		filepath.Join("deployment", "workdir", "main", "user", user, "file1.tf"),
		filepath.Join("deployment", "workdir", "main", "user", user, "file2.tf"),
		filepath.Join("deployment", "workdir", "modules", "module1", "file1.tf"),
		filepath.Join("deployment", "workdir", "modules", "module1", "file2.tf"),
		filepath.Join("deployment", "workdir", "modules", "module2", "file1.tf"),
		filepath.Join("deployment", "workdir", "modules", "module2", "file2.tf"),
	}
}

func testWorkdirCTDDirectories(path string, filePrefix string) []string {
	return []string{
		filepath.Join(path, "main"),
		filepath.Join(path, "main", "global"),
		filepath.Join(path, "main", "user"),
		filepath.Join(path, "main", "user", "user1"),
		filepath.Join(path, "main", "user", "user2"),
		filepath.Join(path, "modules"),
		filepath.Join(path, "modules", filePrefix+"module1"),
		filepath.Join(path, "modules", filePrefix+"module2"),
	}
}

func testWorkdirCTDFiles(path string, filePrefix string) []string {
	return []string{
		filepath.Join(path, "main", "global", filePrefix+"file1.tf"),
		filepath.Join(path, "main", "global", filePrefix+"file2.tf"),
		filepath.Join(path, "main", "user", "user1", filePrefix+"file1.tf"),
		filepath.Join(path, "main", "user", "user1", filePrefix+"file2.tf"),
		filepath.Join(path, "main", "user", "user2", filePrefix+"file1.tf"),
		filepath.Join(path, "main", "user", "user2", filePrefix+"file2.tf"),
		filepath.Join(path, "modules", filePrefix+"module1", "file1.tf"),
		filepath.Join(path, "modules", filePrefix+"module1", "file2.tf"),
		filepath.Join(path, "modules", filePrefix+"module2", "file1.tf"),
		filepath.Join(path, "modules", filePrefix+"module2", "file2.tf"),
	}
}

func testWorktreeCreateCTDFiles(fs afero.Fs, path string, filePrefix string) error {
	for _, directory := range testWorkdirCTDDirectories(path, filePrefix) {
		err := fs.MkdirAll(directory, 0700)
		if err != nil {
			return err
		}
	}

	for _, file := range testWorkdirCTDFiles(path, filePrefix) {
		err := utils.NewFileWithContentIfNotExist(fs, file, file)
		if err != nil {
			return err
		}
	}

	return nil
}
