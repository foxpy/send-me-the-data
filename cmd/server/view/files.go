package view

import (
	"fmt"

	"github.com/foxpy/send-me-the-data/cmd/server/idb"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs"
	"github.com/foxpy/send-me-the-data/cmd/server/template"
)

func Files(link idb.Link, files []ifs.File, offset, limit int) []template.FileView {
	if offset > len(files) {
		return nil
	}
	files = files[offset:]

	if limit <= len(files) {
		files = files[:limit]
	}

	userDownloadable := link.UserDownloadable()
	linkID := link.ID()

	fileViews := make([]template.FileView, 0, len(files))
	for _, file := range files {
		userDownloadLink := ""
		if userDownloadable {
			userDownloadLink = fmt.Sprintf("/%s/%s", linkID, file.Name)
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

	return fileViews
}
