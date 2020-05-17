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

type binary struct {
	fs   afero.Fs
	path string

	version string
	arch    string
}

// BinaryPath returns the path of the terraform binary
func (b *binary) BinaryPath() string {
	return filepath.Join(b.path, fmt.Sprintf("terraform_%s_%s", b.version, b.arch))
}

// getBinary downloads a terraform binary from Hashicorp official release page.
func (b *binary) getBinary() error {
	url := b.terraformDownloadURL()

	zipFile, err := b.downloadZip(url)
	if err != nil {
		return err
	}

	err = b.uncompressZip(zipFile, b.BinaryPath())
	if err != nil {
		return err
	}

	return nil
}

// XXX: a sha check would be great. For now, only check if binary exists
func (b *binary) checkBinary() (bool, error) {
	return afero.Exists(b.fs, b.BinaryPath())
}

func (b *binary) downloadZip(url string) (string, error) {
	filePath := ""
	file, err := afero.TempFile(b.fs, b.path, "terraform_zip_")
	if err != nil {
		return filePath, err
	}
	defer file.Close()

	filePath = file.Name()
	err = utils.HTTPDownloadFile(b.fs, file, url)
	if err != nil {
		return filePath, err
	}

	return filePath, nil
}

func (b *binary) uncompressZip(zipFilePath string, binaryFilePath string) error {
	zipFile, err := b.fs.Open(zipFilePath)
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

	binaryFile, err := b.fs.OpenFile(binaryFilePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0755)
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

func (b *binary) terraformDownloadURL() string {
	return fmt.Sprintf("https://releases.hashicorp.com/terraform/%s/terraform_%s_%s.zip", b.version, b.version, b.arch)
}
