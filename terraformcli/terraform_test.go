package terraformcli

import (
	"crypto/sha256"
	"encoding/hex"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/afero"
)

// This test requires internet access and actually downloads a fully terraform zip.
// We should avoid this.
func TestGetBinary(t *testing.T) {
	// Hardcoded test values
	version := "0.12.24"
	arch := "darwin_amd64"
	expectedSHA := "dfce338efc62080ad02b14c3d389db0d8e33664994373f840ba4001b1c860392"

	fs := afero.NewMemMapFs()

	terraform, err := New(fs, filepath.Join("terraform"), version, arch)
	if err != nil {
		t.Fatal(err)
	}

	err = terraform.GetBinary()
	if err != nil {
		t.Fatal(err)
	}

	terraformBinary, err := afero.ReadFile(fs, terraform.BinaryPath())
	if err != nil {
		t.Fatal(err)
	}
	obtainedSHAbytes := sha256.Sum256(terraformBinary)
	obtainedSHA := hex.EncodeToString(obtainedSHAbytes[:])

	if !reflect.DeepEqual(expectedSHA, obtainedSHA) {
		t.Errorf("Incorrect file list.\n\n Expected: %v\n\n Obtained: %v\n", expectedSHA, obtainedSHA)
	}
}
