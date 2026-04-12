package view

import (
	"fmt"

	"github.com/foxpy/send-me-the-data/cmd/server/idb"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs"
	"github.com/foxpy/send-me-the-data/cmd/server/template"
)

func Files(fs ifs.Filesystem, lock idb.LinkRLock) ([]template.FileView, error) {
	linkID := lock.ExternalKey()
	files, err := fs.ListLinkFiles(linkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all files for link %s: %w", linkID, err)
	}

	fileViews := make([]template.FileView, 0, len(files))
	for _, file := range files {
		userDownloadLink := ""
		if lock.UserDownloadable() {
			userDownloadLink = fmt.Sprintf("/u/%s/%s", linkID, file.Name)
		}

		fileViews = append(fileViews, template.FileView{
			Name:              file.Name,
			UploadedAt:        uint64(file.ModTime.UnixMilli()),
			Size:              bytesToHuman(uint64(file.Size)),
			AdminDownloadLink: fmt.Sprintf("/link/%s/file/%s", linkID, file.Name),
			UserDownloadLink:  userDownloadLink,
			DeleteLink:        fmt.Sprintf("/link/%s/file/%s/delete", linkID, file.Name),
		})
	}

	return fileViews, nil
}
