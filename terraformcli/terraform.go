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

	version string
	arch    string
}

// New constructs a new Terraform struct and returns it
func New(fs afero.Fs, path string, version string, arch string) (Terraform, error) {
	terraform := Terraform{
		fs:   fs,
		path: path,

		version: version,
		arch:    arch,
	}

	err := fs.MkdirAll(path, 0755)
	if err != nil {
		return terraform, err
	}

	ok, err := terraform.checkBinary()
	if err != nil {
		return terraform, err
	}
	if !ok {
		err = terraform.GetBinary()
		if err != nil {
			return terraform, err
		}
	}

	return terraform, nil
}

// BinaryPath returns the path of the terraform binary
func (t *Terraform) BinaryPath() string {
	return filepath.Join(t.path, fmt.Sprintf("terraform_%s_%s", t.version, t.arch))
}

// GetBinary downloads a terraform binary from Hashicorp official release page.
func (t *Terraform) GetBinary() error {
	url := t.terraformDownloadURL()

	zipFile, err := t.downloadZip(url)
	if err != nil {
		return err
	}

	err = t.uncompressZip(zipFile, t.BinaryPath())
	if err != nil {
		return err
	}

	return nil
}

// XXX: a sha check would be great. For now, only check if binary exists
func (t *Terraform) checkBinary() (bool, error) {
	return afero.Exists(t.fs, t.BinaryPath())
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

func (t *Terraform) terraformDownloadURL() string {
	return fmt.Sprintf("https://releases.hashicorp.com/terraform/%s/terraform_%s_%s.zip", t.version, t.version, t.arch)
}
