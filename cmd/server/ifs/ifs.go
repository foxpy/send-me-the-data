package ifs

import (
	"io/fs"
	"os"
	"time"
)

type Filesystem interface {
	ListLinkFiles(linkID string) ([]File, error)
	RemoveLinkFiles(linkID string) error
	RemoveLinkFile(linkID, fileName string) error
	LinkFS(linkID string) (fs.FS, error)
	CreateNewFile(linkID, fileName string) (*os.File, error)
}

type File struct {
	Name string
	Size int64
	// A wiser choice would be to rely on creation time instead,
	// but since it is not required by POSIX, it is not supported by Go either
	ModTime time.Time
}
