package deployment

import (
	"reflect"
	"testing"

	"github.com/arodriguezdlc/sonatina/utils"

	"github.com/spf13/afero"
)

func TestListMainGlobalFiles(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := testCreateFsTree(fs)
	if err != nil {
		t.Error(err)
	}

	ctd, err := NewCTD(fs, "/", "example.com", "/")
	if err != nil {
		t.Error(err)
	}

	obtainedFileList, err := ctd.ListMainGlobalFiles()
	if err != nil {
		t.Error(err)
	}
	expectedFileList := []string{
		"/main/global/file1.tf",
		"/main/global/file2.tf",
	}

	if !reflect.DeepEqual(expectedFileList, obtainedFileList) {
		t.Errorf("Incorrect file list.\n\n Expected: %v\n\n Obtained: %v\n", expectedFileList, obtainedFileList)
	}
}

func TestListUserFiles(t *testing.T) {
	fs := afero.NewMemMapFs()

	err := testCreateFsTree(fs)
	if err != nil {
		t.Error(err)
	}

	ctd, err := NewCTD(fs, "/", "example.com", "/")
	if err != nil {
		t.Error(err)
	}

	obtainedFileList, err := ctd.ListMainUserFiles("user1")
	if err != nil {
		t.Error(err)
	}
	expectedFileList := []string{
		"/main/user/user1/file1.tf",
		"/main/user/user1/file2.tf",
	}

	if !reflect.DeepEqual(expectedFileList, obtainedFileList) {
		t.Errorf("Incorrect file list.\n\n Expected: %v\n\n Obtained: %v\n", expectedFileList, obtainedFileList)
	}

	obtainedFileList, err = ctd.ListMainUserFiles("user2")
	if err != nil {
		t.Error(err)
	}
	expectedFileList = []string{
		"/main/user/user2/file1.tf",
		"/main/user/user2/file2.tf",
	}

	if !reflect.DeepEqual(expectedFileList, obtainedFileList) {
		t.Errorf("Incorrect file list.\n\n Expected: %v\n\n Obtained: %v\n", expectedFileList, obtainedFileList)
	}
}

func testCreateFsTree(fs afero.Fs) error {
	for _, directory := range testCTDDirectories() {
		err := fs.Mkdir(directory, 0700)
		if err != nil {
			return err
		}
	}

	for _, file := range testCTDFiles() {
		err := utils.NewFileIfNotExist(fs, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func testCTDDirectories() []string {
	return []string{
		"/main",
		"/main/global",
		"/main/user",
		"/main/user/user1",
		"/main/user/user2",
		"/modules/module1",
		"/modules/module2",
		"/vtd/",
	}
}

func testCTDFiles() []string {
	return []string{
		"/main/global/file1.tf",
		"/main/global/file2.tf",
		"/main/user/user1/file1.tf",
		"/main/user/user1/file2.tf",
		"/main/user/user2/file1.tf",
		"/main/user/user2/file2.tf",
		"/modules/module1/file1.tf",
		"/modules/module1/file2.tf",
		"/modules/module2/file1.tf",
		"/modules/module2/file2.tf",
	}
}
