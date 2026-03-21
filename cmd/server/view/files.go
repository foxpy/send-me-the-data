package view

import (
	"fmt"

	"github.com/foxpy/send-me-the-data/cmd/server/ifs"
	"github.com/foxpy/send-me-the-data/cmd/server/template"
)

func Files(fs ifs.Filesystem, linkID string, renderDownloadLinks bool) ([]template.FileView, error) {
	files, err := fs.ListLinkFiles(linkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all files for link %s: %w", linkID, err)
	}

	fileViews := make([]template.FileView, 0, len(files))
	for _, file := range files {
		downloadLink := ""
		if renderDownloadLinks {
			downloadLink = fmt.Sprintf("/link/%s/file/%s", linkID, file.Name)
		}

		fileViews = append(fileViews, template.FileView{
			Name:         file.Name,
			UploadedAt:   file.ModTime.Format(DateTimeFormat),
			Size:         bytesToHuman(file.Size),
			DownloadLink: downloadLink,
			DeleteLink:   fmt.Sprintf("/link/%s/file/%s/delete", linkID, file.Name),
		})
	}

	return fileViews, nil
}
