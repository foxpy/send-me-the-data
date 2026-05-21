package mockfs

import (
	"io/fs"
	"os"

	"github.com/foxpy/send-me-the-data/src/ifs"
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
	panic("not implemented")
}

func (f *MockFS) RemoveLinkFile(linkID, fileName string) error {
	panic("not implemented")
}

func (f *MockFS) LinkFS(linkID string) (fs.FS, error) {
	panic("not implemented")
}

func (f *MockFS) CreateNewFile(linkID, fileName string) (*os.File, error) {
	panic("not implemented")
}
