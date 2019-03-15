package utils

import (
	"os"

	"github.com/spf13/afero"
)

// NewFileIfNotExist creates a empty file if is not present, or do
// nothing if file already exists
func NewFileIfNotExist(path string, fs afero.Fs) error {
	var err error
	var file afero.File

	if file, err = fs.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		return err
	}
	defer file.Close()
	return nil
}
