package utils

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/afero"
)

// HTTPDownloadFile downloads a url and store it in local filepath.
// Based on https://progolang.com/how-to-download-files-in-go/
func HTTPDownloadFile(fs afero.Fs, file afero.File, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = httpCheckCorrectResponseCode(resp)
	if err != nil {
		return err
	}

	// Write the body to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func httpCheckCorrectResponseCode(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return HTTPBadResponseError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
		}
	}
	return nil
}

type HTTPBadResponseError struct {
	StatusCode int
	Status     string
}

func (err HTTPBadResponseError) Error() string {
	return fmt.Sprintf("HTTP Response with Status Code %s", err.Status)
}
