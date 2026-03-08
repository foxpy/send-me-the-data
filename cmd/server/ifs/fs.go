package ifs

import (
	"io/fs"
	"os"
	"time"
)

type Filesystem interface {
	// FIXME: do not read all dir entries into memory at once, use pagination instead
	ListLinkFiles(linkID string) ([]File, error)
	RemoveLinkFiles(linkID string) error
	RemoveLinkFile(linkID, fileName string) error
	FS(linkID string) (fs.FS, error)
	CreateNewFile(linkID, fileName string) (*os.File, error)
}

type File struct {
	Name    string
	Size    int64
	ModTime time.Time
}
