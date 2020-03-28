package utils

import (
	"os"

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

	ok, err := FileHasContent(file)
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

func FileHasContent(file afero.File) (bool, error) {
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
