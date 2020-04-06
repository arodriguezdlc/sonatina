package utils

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/afero"
)

// NewFileIfNotExist creates a empty file if is not present, or do
// nothing if file already exists
func NewFileIfNotExist(fs afero.Fs, path string) error {
	file, err := fs.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewFileWithContentIfNotExist(fs afero.Fs, path string, content string) error {
	file, err := fs.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return err
	}

	ok, err := fileHasContent(file)
	if err != nil {
		return err
	}
	if !ok {
		_, err := file.WriteString(content)
		if err != nil {
			return err
		}
	}

	return nil
}

func FileCopy(fs afero.Fs, sourcePath string, destPath string) error {
	// XXX: easy way to copy file, it isn't good for large files because has to save it
	// completely in memory before write. However, this frameworks usually works with small files.
	// It could be improved in the future.
	fileInfo, err := fs.Stat(sourcePath)
	if err != nil {
		return err
	}

	fileContent, err := afero.ReadFile(fs, sourcePath)
	if err != nil {
		return err
	}

	return afero.WriteFile(fs, destPath, fileContent, fileInfo.Mode())
}

// Based on walk function https://github.com/spf13/afero/blob/master/path.go
func FileCopyRecursively(fs afero.Fs, sourcePath string, destPath string) error {
	// TODO
	filenames, err := readDirNames(fs, sourcePath)
	if err != nil {
		return err
	}
	for _, filename := range filenames {
		sourceFilePath := filepath.Join(sourcePath, filename)
		destFilePath := filepath.Join(destPath, filename)

		sourceFileInfo, err := fs.Stat(sourceFilePath)
		if err != nil {
			return err
		}

		if sourceFileInfo.IsDir() {
			err = fs.Mkdir(destFilePath, sourceFileInfo.Mode())
			if err != nil {
				return err
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

func FileListRecursively(fs afero.Fs, path string) ([]string, error) {
	files, err := readDirNamesGlob(fs, path)
	if err != nil {
		return files, err
	}

	for _, file := range files {
		fileInfo, err := fs.Stat(file)
		if err != nil {
			return files, err
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

func fileHasContent(file afero.File) (bool, error) {
	result := false

	fileInfo, err := file.Stat()
	if err != nil {
		return result, err
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
		return nil, err
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	sort.Strings(names)
	return names, nil
}

// Adapted from https://github.com/spf13/afero/blob/master/path.go
func readDirNamesGlob(fs afero.Fs, dirname string) ([]string, error) {
	f, err := fs.Open(dirname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	namesGlob := []string{}
	for _, name := range names {
		nameGlob := filepath.Join(dirname, name)
		namesGlob = append(namesGlob, nameGlob)
	}

	sort.Strings(namesGlob)
	return namesGlob, nil
}
