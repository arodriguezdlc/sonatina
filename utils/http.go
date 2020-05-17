package utils

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// HTTPDownloadFile downloads a url and store it in local filepath.
// Based on https://progolang.com/how-to-download-files-in-go/
func HTTPDownloadFile(fs afero.Fs, file afero.File, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrapf(err, "couldn't perform HTTP GET over %s", url)
	}
	defer resp.Body.Close()

	err = httpCheckCorrectResponseCode(resp)
	if err != nil {
		return err
	}

	// Write the body to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return errors.Wrap(err, "couldn't copy streams")
	}
	return nil
}

func httpCheckCorrectResponseCode(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(fmt.Sprintf("HTTP Response with Status Code %s", resp.Status))
	}
	return nil
}
