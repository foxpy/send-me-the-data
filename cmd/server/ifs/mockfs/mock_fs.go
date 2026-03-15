package mockfs

import (
	"io/fs"
	"os"

	"github.com/foxpy/send-me-the-data/cmd/server/ifs"
)

type MockFS struct {
	listLinkFilesResponses map[string][]ifs.File
}

var _ ifs.Filesystem = &MockFS{}

func NewMockFS() *MockFS {
	return &MockFS{
		listLinkFilesResponses: make(map[string][]ifs.File),
	}
}

func (f *MockFS) RemoveLinkFiles(linkID string) error {
	panic("TODO")
}

func (f *MockFS) RemoveLinkFile(linkID, fileName string) error {
	panic("TODO")
}

func (f *MockFS) FS(linkID string) (fs.FS, error) {
	panic("TODO")
}

func (f *MockFS) CreateNewFile(linkID, fileName string) (*os.File, error) {
	panic("TODO")
}
