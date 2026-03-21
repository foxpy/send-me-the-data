package vfs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/foxpy/send-me-the-data/cmd/server/ifs"
)

type VFS struct {
	root *os.Root
}

var _ ifs.Filesystem = &VFS{}

func NewVFS(prefix string) (*VFS, error) {
	if prefix == "" {
		return nil, errors.New("filesystem prefix must not be an empty string")
	}

	root, err := os.OpenRoot(prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to open filesystem root at %s: %w", prefix, err)
	}

	return &VFS{root}, nil
}

func (f *VFS) ListLinkFiles(linkID string) ([]ifs.File, error) {
	linkFolder, err := f.root.Open(linkID)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to open directory of link %s: %w", linkID, err)
	}

	defer func() {
		_ = linkFolder.Close()
	}()

	entries, err := linkFolder.ReadDir(0)
	if err != nil {
		return nil, fmt.Errorf("failed to traverse directory of link %s: %w", linkID, err)
	}

	files := make([]ifs.File, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("failed to get file info: %w", err)
		}

		files = append(files, ifs.File{
			Name:    entry.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime().UTC(),
		})
	}

	return files, nil
}

func (f *VFS) RemoveLinkFiles(linkID string) error {
	err := f.root.RemoveAll(linkID)
	if err != nil {
		return fmt.Errorf("failed to remove files associated with link %s: %w", linkID, err)
	}

	return nil
}

func (f *VFS) RemoveLinkFile(linkID, fileName string) error {
	return f.root.Remove(f.getPath(linkID, fileName))
}

func (f *VFS) FS(linkID string) (fs.FS, error) {
	linkRoot, err := f.root.OpenRoot(linkID)
	if err != nil {
		return nil, err
	}

	return linkRoot.FS(), err
}

func (f *VFS) CreateNewFile(linkID, fileName string) (*os.File, error) {
	err := f.root.MkdirAll(linkID, 0777)
	if err != nil {
		return nil, fmt.Errorf("failed to create a directory for a new file %s for a link %s: %w", fileName, linkID, err)
	}

	return f.root.OpenFile(f.getPath(linkID, fileName), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
}

func (f *VFS) getPath(linkID, fileName string) string {
	return fmt.Sprintf("%s/%s", linkID, fileName)
}
