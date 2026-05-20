package view

import (
	"fmt"

	"github.com/foxpy/send-me-the-data/src/idb"
	"github.com/foxpy/send-me-the-data/src/ifs"
	"github.com/foxpy/send-me-the-data/src/template"
)

func Link(link idb.Link, files []ifs.File) template.LinkView {
	var totalSize uint64
	for _, file := range files {
		totalSize += uint64(file.Size)
	}

	return template.LinkView{
		Name:             link.Name(),
		CreatedAt:        uint64(link.CreatedAt().UnixMilli()),
		TotalFiles:       len(files),
		TotalSize:        bytesToHuman(totalSize),
		MaxFileSize:      bytesToHuman(link.MaxFileSize()),
		MaxFileSizeBytes: link.MaxFileSize(),
		ViewLink:         fmt.Sprintf("/link/%s", link.ID()),
		DeleteLink:       fmt.Sprintf("/link/%s/delete", link.ID()),
		EditLink:         fmt.Sprintf("/link/%s/edit", link.ID()),
		DownloadZIP:      fmt.Sprintf("/link/%s/zip", link.ID()),
		UserDownloadable: link.UserDownloadable(),
		UploadEnabled:    link.UploadEnabled(),
	}
}
