package terraformcli

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/arodriguezdlc/sonatina/utils"
	"github.com/spf13/afero"
)

type Terraform struct {
	fs   afero.Fs
	path string
}

// New constructs a new Terraform struct and returns it
func New(fs afero.Fs, path string) (Terraform, error) {
	terraform := Terraform{
		fs:   fs,
		path: path,
	}

	err := fs.MkdirAll(path, 0755)
	if err != nil {
		return terraform, err
	}

	return terraform, nil
}

// BinaryPath returns the path of the terraform binary with a specified version and arch.
func (t *Terraform) BinaryPath(version string, arch string) string {
	return filepath.Join(t.path, fmt.Sprintf("terraform_%s_%s", version, arch))
}

// GetBinary downloads a terraform binary from Hashicorp official release page.
func (t *Terraform) GetBinary(version string, arch string) error {
	url := t.terraformDownloadURL(version, arch)

	zipFile, err := t.downloadZip(url)
	if err != nil {
		return err
	}

	err = t.uncompressZip(zipFile, t.BinaryPath(version, arch))
	if err != nil {
		return err
	}

	return nil
}

func (t *Terraform) downloadZip(url string) (string, error) {
	filePath := ""
	file, err := afero.TempFile(t.fs, t.path, "terraform_zip_")
	if err != nil {
		return filePath, err
	}
	defer file.Close()

	filePath = file.Name()
	err = utils.HTTPDownloadFile(t.fs, file, url)
	if err != nil {
		return filePath, err
	}

	return filePath, nil
}

func (t *Terraform) uncompressZip(zipFilePath string, binaryFilePath string) error {
	zipFile, err := t.fs.Open(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipFileInfo, err := zipFile.Stat()
	if err != nil {
		return err
	}

	reader, err := zip.NewReader(zipFile, zipFileInfo.Size())
	if err != nil {
		return err
	}

	binaryFile, err := t.fs.OpenFile(binaryFilePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer binaryFile.Close()

	file, err := reader.File[0].Open()
	if err != nil {
		return err
	}

	_, err = io.Copy(binaryFile, file)
	if err != nil {
		return err
	}

	return nil
}

func (t *Terraform) terraformDownloadURL(version string, arch string) string {
	return fmt.Sprintf("https://releases.hashicorp.com/terraform/%s/terraform_%s_%s.zip", version, version, arch)
}
