package deployment

import (
	"testing"

	"github.com/spf13/afero"
)

func TestListMainFilesGlobal(t *testing.T) {
	fs := afero.NewMemMapFs()

	ctd, err := NewCTD(fs, "/", "example.com", "/")
	if err != nil {
		t.Error(err)
	}

	//obtained_file_list, err := ctd.ListMainFilesGlobal()

}

func testCreateFsTree(fs afero.Fs) error {
	directories := []string{
		"/main",
		"/main/global",
		"/main/user",
		"/main/user/user1",
		"/main/user/user2",
		"/modules/module1",
		"/modules/module2",
		"/vtd/",
	}

	for _, directory := range directories {
		err := fs.Mkdir(directory, 0700)
		if err != nil {
			return err
		}
	}

	return nil
}
