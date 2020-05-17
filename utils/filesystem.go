package utils

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// NewFileIfNotExist creates a empty file if is not present, or do
// nothing if file already exists
func NewFileIfNotExist(fs afero.Fs, path string) error {
	file, err := fs.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return errors.Wrapf(err, "couldn't open file %s", path)
	}
	return nil
}

// NewFileWithContentIfNotExist creates a file with a speficied content only if file
// doesn't exist
func NewFileWithContentIfNotExist(fs afero.Fs, path string, content string) error {
	file, err := fs.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return errors.Wrapf(err, "couldn't open file %s", path)
	}

	ok, err := fileHasContent(file)
	if err != nil {
		return err
	}
	if !ok {
		_, err := file.WriteString(content)
		if err != nil {
			return errors.Wrapf(err, "couldn't write on file %s", path)
		}
	}

	return nil
}

// NewDirectoryWithKeep creates a directory and an empty .keep file to maintain
// filetree in git
func NewDirectoryWithKeep(fs afero.Fs, path string) error {
	err := fs.MkdirAll(path, 0755)
	if err != nil {
		return errors.Wrapf(err, "couldn't create directory %s", path)
	}

	err = NewFileIfNotExist(fs, filepath.Join(path, ".keep"))
	if err != nil {
		return err
	}

	return nil
}

func FileCopy(fs afero.Fs, sourcePath string, destPath string) error {
	// XXX: easy way to copy file, it isn't good for large files because has to save it
	// completely in memory before write. However, this frameworks usually works with small files.
	// It could be improved in the future.
	fileInfo, err := fs.Stat(sourcePath)
	if err != nil {
		return errors.Wrapf(err, "couldn't get stat from file %s", sourcePath)
	}

	fileContent, err := afero.ReadFile(fs, sourcePath)
	if err != nil {
		return errors.Wrapf(err, "couldn't read from file %s", sourcePath)
	}

	err = afero.WriteFile(fs, destPath, fileContent, fileInfo.Mode())
	if err != nil {
		return errors.Wrapf(err, "couldn't write to file %s", destPath)
	}

	return nil
}

// FileCopyRecursively copies recursively files from one directory to other, being
// similar to a `cp -R` command.
// Based on walk function https://github.com/spf13/afero/blob/master/path.go
func FileCopyRecursively(fs afero.Fs, sourcePath string, destPath string) error {
	filenames, err := readDirNames(fs, sourcePath)
	if err != nil {
		return err
	}
	for _, filename := range filenames {
		sourceFilePath := filepath.Join(sourcePath, filename)
		destFilePath := filepath.Join(destPath, filename)

		sourceFileInfo, err := fs.Stat(sourceFilePath)
		if err != nil {
			return errors.Wrapf(err, "couldn't get stat from file %s", sourceFilePath)
		}

		if sourceFileInfo.IsDir() {
			err = fs.Mkdir(destFilePath, sourceFileInfo.Mode())
			if err != nil {
				return errors.Wrapf(err, "couldn't create directory %s", destFilePath)
			}
			err = FileCopyRecursively(fs, sourceFilePath, destFilePath)
			if err != nil {
				return err
			}
		} else {
			err = FileCopy(fs, sourceFilePath, destFilePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// FileListRecursively list recursively all files inside a directory. The result list
// is alphanumerically ordered
func FileListRecursively(fs afero.Fs, path string) ([]string, error) {
	files, err := readDirNamesGlob(fs, path)
	if err != nil {
		return files, err
	}

	for _, file := range files {
		fileInfo, err := fs.Stat(file)
		if err != nil {
			return files, errors.Wrapf(err, "couldn't get stat from file %s", file)
		}
		if fileInfo.IsDir() {
			appendFiles, err := FileListRecursively(fs, file)
			if err != nil {
				return files, err
			}
			files = append(files, appendFiles...)
		}
	}

	sort.Strings(files)
	return files, nil
}

// FileListRecursivelyWithoutDirs list recursively all files inside a directory, excluding
// directories. The result list is alphanumerically ordered
func FileListRecursivelyWithoutDirs(fs afero.Fs, path string) ([]string, error) {
	// This is not an efficient implementation, but is quick and easy to implement. Could be improved.
	filteredFileList := []string{}
	fileList, err := FileListRecursively(fs, path)
	if err != nil {
		return filteredFileList, err
	}

	for _, file := range fileList {
		fileInfo, err := fs.Stat(file)
		if err != nil {
			return filteredFileList, errors.Wrapf(err, "couldn't get stat from file %s", file)
		}
		if !fileInfo.IsDir() {
			filteredFileList = append(filteredFileList, file)
		}
	}

	return filteredFileList, err
}

func fileHasContent(file afero.File) (bool, error) {
	result := false

	fileInfo, err := file.Stat()
	if err != nil {
		return result, errors.Wrapf(err, "couldn't get stat from file %s", file.Name())
	}

	if fileInfo.Size() > 0 {
		result = true
	}
	return result, nil
}

// Adapted from https://github.com/spf13/afero/blob/master/path.go
func readDirNames(fs afero.Fs, dirname string) ([]string, error) {
	f, err := fs.Open(dirname)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't open file %s", dirname)
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't read dir names from %s", dirname)
	}

	sort.Strings(names)
	return names, nil
}

// Adapted from https://github.com/spf13/afero/blob/master/path.go
func readDirNamesGlob(fs afero.Fs, dirname string) ([]string, error) {
	f, err := fs.Open(dirname)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't open file %s", dirname)
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't read dir names from %s", dirname)
	}

	namesGlob := []string{}
	for _, name := range names {
		nameGlob := filepath.Join(dirname, name)
		namesGlob = append(namesGlob, nameGlob)
	}

	sort.Strings(namesGlob)
	return namesGlob, nil
}
