package view

import (
	"fmt"

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

		var totalSize int64
		for _, file := range files {
			totalSize += file.Size
		}

		linkViews = append(linkViews, template.LinkView{
			Name:       link.Name,
			CreatedAt:  link.CreatedAt.Format(DateTimeFormat),
			TotalFiles: len(files),
			TotalSize:  bytesToHuman(totalSize),
			ViewLink:   fmt.Sprintf("/link/%s", link.ExternalKey),
			DeleteLink: fmt.Sprintf("/link/%s/delete", link.ExternalKey),
		})
	}

	return linkViews, nil
}
