package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestHTTPDownloadFile(t *testing.T) {
	expectedContent := "Test"
	filename := "/downloadedFile"

	fs := afero.NewMemMapFs()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, expectedContent)
	}))
	defer ts.Close()

	// Create the file
	file, err := fs.OpenFile(filename, os.O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	err = HTTPDownloadFile(fs, file, ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	obtainedContentBytes, err := afero.ReadFile(fs, filename)
	if err != nil {
		t.Fatal(err)
	}
	obtainedContent := strings.TrimSpace(string(obtainedContentBytes))

	if !reflect.DeepEqual(expectedContent, obtainedContent) {
		t.Errorf("Incorrect downloaded content.\n\n Expected: %v\n\n Obtained: %v\n", expectedContent, obtainedContent)
	}
}
