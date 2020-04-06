package utils

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestNewFileIfNotExistsWithExistingFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	expectedContent := "test file content"
	file := "/file"

	err := NewFileWithContentIfNotExist(fs, file, expectedContent)
	if err != nil {
		t.Error(err)
	}

	err = NewFileIfNotExist(fs, file)
	if err != nil {
		t.Error(err)
	}

	obtainedContentBytes, err := afero.ReadFile(fs, file)
	obtainedContent := string(obtainedContentBytes)
	if !reflect.DeepEqual(expectedContent, obtainedContent) {
		t.Errorf("Incorrect file content.\n\n Expected: %v\n\n Obtained: %v\n", expectedContent, obtainedContent)
	}
}

func TestNewFileIfNotExistsWithUnexistingFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	file := "/file"
	expectedContent := ""

	err := NewFileIfNotExist(fs, file)
	if err != nil {
		t.Error(err)
	}

	obtainedContentBytes, err := afero.ReadFile(fs, file)
	obtainedContent := string(obtainedContentBytes)
	if !reflect.DeepEqual(expectedContent, obtainedContent) {
		t.Errorf("Incorrect file content.\n\n Expected empty file\n\n Obtained: %v\n", obtainedContent)
	}
}

func TestNewFileWithContentWithExistingFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	file := "/file"
	expectedContent := "test file content"

	err := NewFileWithContentIfNotExist(fs, file, expectedContent)
	if err != nil {
		t.Error(err)
	}

	err = NewFileWithContentIfNotExist(fs, file, "other content")
	if err != nil {
		t.Error(err)
	}

	obtainedContentBytes, err := afero.ReadFile(fs, file)
	obtainedContent := string(obtainedContentBytes)
	if !reflect.DeepEqual(expectedContent, obtainedContent) {
		t.Errorf("Incorrect file content.\n\n Expected empty file\n\n Obtained: %v\n", obtainedContent)
	}
}

func TestNewFileWithContentWithUnexistingFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	file := "/file"
	expectedContent := "test file content"

	err := NewFileWithContentIfNotExist(fs, file, expectedContent)
	if err != nil {
		t.Error(err)
	}

	obtainedContentBytes, err := afero.ReadFile(fs, file)
	obtainedContent := string(obtainedContentBytes)
	if !reflect.DeepEqual(expectedContent, obtainedContent) {
		t.Errorf("Incorrect file content.\n\n Expected empty file\n\n Obtained: %v\n", obtainedContent)
	}
}

func TestFileCopy(t *testing.T) {
	fs := afero.NewMemMapFs()
	expectedContent := "test file content"
	sourceFile := "/source_file"
	destFile := "/dest_file"

	err := NewFileWithContentIfNotExist(fs, sourceFile, expectedContent)
	if err != nil {
		t.Error(err)
	}

	FileCopy(fs, sourceFile, destFile)

	obtainedContentBytes, err := afero.ReadFile(fs, destFile)
	obtainedContent := string(obtainedContentBytes)
	if !reflect.DeepEqual(expectedContent, obtainedContent) {
		t.Errorf("Incorrect file content.\n\n Expected: %v\n\n Obtained: %v\n", expectedContent, obtainedContent)
	}

	sourceStat, err := fs.Stat(sourceFile)
	if err != nil {
		t.Error(err)
	}
	sourceMode := sourceStat.Mode()
	destStat, err := fs.Stat(destFile)
	if err != nil {
		t.Error(err)
	}
	destMode := destStat.Mode()

	if !reflect.DeepEqual(sourceMode, destMode) {
		t.Errorf("Incorrect file mode.\n\n Expected: %v\n\n Obtained: %v\n", sourceMode, destMode)
	}
}

func TestFileCopyRecursively(t *testing.T) {
	fs := afero.NewMemMapFs()
	sourcePath := "/source"
	destPath := "/dest"
	expectedFileTree := testFileTreeTrimmed(destPath)

	err := testCreateFileTree(fs, sourcePath)
	if err != nil {
		t.Error(err)
	}

	err = FileCopyRecursively(fs, sourcePath, destPath)
	if err != nil {
		t.Error(err)
	}

	obtainedFileTree, err := FileListRecursively(fs, "/dest")
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expectedFileTree, obtainedFileTree) {
		t.Errorf("Incorrect file tree.\n\n Expected: %v\n\n Obtained: %v\n", expectedFileTree, obtainedFileTree)
	}

	// Check if file content is correct. Each test file must have a string with
	// its path.
	for _, filePath := range testFiles(sourcePath) {
		expectedContent := filePath
		obtainedContentBytes, err := afero.ReadFile(fs, filePath)
		if err != nil {
			t.Error(err)
		}
		obtainedContent := string(obtainedContentBytes)

		if !reflect.DeepEqual(expectedContent, obtainedContent) {
			t.Errorf("Incorrect content for file %v.\n\n Expected: %v\n\n Obtained: %v\n", filePath, expectedContent, obtainedContent)
		}
	}
}

func TestFileListRecursively(t *testing.T) {
	fs := afero.NewMemMapFs()
	path := "/list"
	expectedFileTree := testFileTreeTrimmed(path)

	err := testCreateFileTree(fs, path)
	if err != nil {
		t.Error(err)
	}

	obtainedFileTree, err := FileListRecursively(fs, path)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expectedFileTree, obtainedFileTree) {
		t.Errorf("Incorrect file tree.\n\n Expected: %v\n\n Obtained: %v\n", expectedFileTree, obtainedFileTree)
	}
}

func testFileTree(basepath string) []string {
	filetree := []string{
		basepath + "/file1.txt",
		basepath + "/folder1/",
		basepath + "/folder2/",
		basepath + "/folder2/file1.txt",
		basepath + "/folder2/file2.txt",
		basepath + "/folder2/file3.txt",
		basepath + "/folder3/",
		basepath + "/folder3/file1.txt",
		basepath + "/folder3/folder1",
		basepath + "/folder3/folder2",
		basepath + "/folder3/folder1/file1.txt",
		basepath + "/folder3/folder1/file2.txt",
		basepath + "/folder3/folder2/file1.txt",
		basepath + "/folder3/folder2/file2.txt",
	}
	sort.Strings(filetree)
	return filetree
}

func testFileTreeTrimmed(basepath string) []string {
	trimmed := []string{}
	for _, file := range testFileTree(basepath) {
		trimmed = append(trimmed, strings.TrimSuffix(file, "/"))
	}
	return trimmed
}

func testCreateFileTree(fs afero.Fs, basepath string) error {
	for _, filePath := range testFileTree(basepath) {
		if strings.HasSuffix(filePath, "/") {
			err := fs.Mkdir(filePath, 0700)
			if err != nil {
				return err
			}
		} else {
			err := NewFileWithContentIfNotExist(fs, filePath, filePath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func testFiles(basepath string) []string {
	files := []string{}
	for _, filePath := range testFileTree(basepath) {
		if !strings.HasSuffix(filePath, "/") {
			files = append(files, filePath)
		}
	}
	return files
}
