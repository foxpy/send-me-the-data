package filesystem

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"
)

type Filesystem struct {
	root *os.Root
}

type File struct {
	Name    string
	Size    int64
	ModTime time.Time
}

func NewFilesystem(prefix string) (*Filesystem, error) {
	if prefix == "" {
		return nil, errors.New("filesystem prefix must not be an empty string")
	}

	root, err := os.OpenRoot(prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to open filesystem root at %s: %w", prefix, err)
	}

	return &Filesystem{root}, nil
}

// TODO: this function should really just accept linkID and fileName
func (f *Filesystem) Remove(path string) error {
	return f.root.Remove(path)
}

func (f *Filesystem) ListLinkFiles(linkID string) ([]File, error) {
	linkFolder, err := f.root.Open(linkID)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to open directory of link %s: %w", linkID, err)
	}

	defer func() {
		_ = linkFolder.Close()
	}()

	// FIXME: do not read all dir entries into memory at once
	// solution: pagination
	entries, err := linkFolder.ReadDir(0)
	if err != nil {
		return nil, fmt.Errorf("failed to traverse directory of link %s: %w", linkID, err)
	}

	files := make([]File, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("failed to get file info: %w", err)
		}

		files = append(files, File{
			Name:    entry.Name(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})
	}

	return files, nil
}

func (f *Filesystem) RemoveLinkFiles(linkID string) error {
	err := f.root.RemoveAll(linkID)
	if err != nil {
		return fmt.Errorf("failed to remove files associated with link %s: %w", linkID, err)
	}

	return nil
}

func (f *Filesystem) RemoveLinkFile(linkID, fileName string) error {
	return f.root.Remove(f.GetPath(linkID, fileName))
}

func (f *Filesystem) FS(linkID string) (fs.FS, error) {
	linkRoot, err := f.root.OpenRoot(linkID)
	if err != nil {
		return nil, err
	}

	return linkRoot.FS(), err
}

// TODO: should this method be inlined into CreateNewFile?
func (f *Filesystem) CreateLinkDirectory(linkID string) error {
	return f.root.MkdirAll(linkID, 0777)
}

func (f *Filesystem) GetPath(linkID, fileName string) string {
	return fmt.Sprintf("%s/%s", linkID, fileName)
}

func (f *Filesystem) CreateNewFile(linkID, fileName string) (*os.File, error) {
	err := f.CreateLinkDirectory(linkID)
	if err != nil {
		return nil, fmt.Errorf("failed to create a directory for a new file %s for a link %s: %w", fileName, linkID, err)
	}

	return f.root.OpenFile(f.GetPath(linkID, fileName), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
}
