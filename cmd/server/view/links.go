package view

import (
	"fmt"
	"strconv"

	"github.com/foxpy/send-me-the-data/cmd/server/idb"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs"
	"github.com/foxpy/send-me-the-data/cmd/server/template"
)

func Links(db idb.Database, fs ifs.Filesystem) ([]template.LinkView, error) {
	links, err := db.AllLinks()
	if err != nil {
		return nil, fmt.Errorf("failed to read all links from database: %w", err)
	}

	linkViews := make([]template.LinkView, 0, len(links))
	for _, link := range links {
		files, err := fs.ListLinkFiles(link.ExternalKey)
		if err != nil {
			return nil, fmt.Errorf("failed to get all files for link %s: %w", link.ExternalKey, err)
		}

		var totalSize uint64
		for _, file := range files {
			totalSize += uint64(file.Size)
		}

		linkViews = append(linkViews, template.LinkView{
			Name:             link.Name,
			CreatedAt:        uint64(link.CreatedAt.UnixMilli()),
			TotalFiles:       len(files),
			TotalSize:        bytesToHuman(totalSize),
			MaxFileSize:      bytesToHuman(link.MaxFileSize),
			MaxFileSizeBytes: strconv.Itoa(int(link.MaxFileSize)),
			ViewLink:         fmt.Sprintf("/link/%s", link.ExternalKey),
			DeleteLink:       fmt.Sprintf("/link/%s/delete", link.ExternalKey),
			EditLink:         fmt.Sprintf("/link/%s/edit", link.ExternalKey),
			DownloadZIP:      fmt.Sprintf("/link/%s/zip", link.ExternalKey),
			UserDownloadable: link.UserDownloadable,
		})
	}

	return linkViews, nil
}

func Link(linkLock idb.LinkRLock, fs ifs.Filesystem) (*template.LinkView, error) {
	id := linkLock.ExternalKey()
	files, err := fs.ListLinkFiles(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get all files for link %s: %w", id, err)
	}

	var totalSize uint64
	for _, file := range files {
		totalSize += uint64(file.Size)
	}

	return &template.LinkView{
		Name:             linkLock.Name(),
		CreatedAt:        uint64(linkLock.CreatedAt().UnixMilli()),
		TotalFiles:       len(files),
		TotalSize:        bytesToHuman(totalSize),
		MaxFileSize:      bytesToHuman(linkLock.MaxFileSize()),
		MaxFileSizeBytes: strconv.Itoa(int(linkLock.MaxFileSize())),
		ViewLink:         fmt.Sprintf("/link/%s", id),
		DeleteLink:       fmt.Sprintf("/link/%s/delete", id),
		EditLink:         fmt.Sprintf("/link/%s/edit", id),
		DownloadZIP:      fmt.Sprintf("/link/%s/zip", id),
		UserDownloadable: linkLock.UserDownloadable(),
	}, nil
}
