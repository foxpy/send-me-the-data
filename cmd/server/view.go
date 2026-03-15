package main

import (
	"fmt"
	"time"

	"github.com/foxpy/send-me-the-data/cmd/server/templates"
)

func (s *State) GetLinksView() ([]templates.LinkView, error) {
	links, err := s.db.AllLinks()
	if err != nil {
		return nil, fmt.Errorf("failed to read all links from database: %w", err)
	}

	linkViews := make([]templates.LinkView, 0, len(links))
	for _, link := range links {
		files, err := s.fs.ListLinkFiles(link.ExternalKey)
		if err != nil {
			return nil, fmt.Errorf("failed to get all files for link %s: %w", link.ExternalKey, err)
		}

		var totalSize int64
		for _, file := range files {
			totalSize += file.Size
		}

		linkViews = append(linkViews, templates.LinkView{
			Name: link.Name,
			// FIXME: using time.Stamp does not include the year
			CreatedAt:  link.CreatedAt.Format(time.Stamp),
			TotalFiles: len(files),
			TotalSize:  bytesToHuman(totalSize),
			ViewLink:   fmt.Sprintf("/link/%s", link.ExternalKey),
			DeleteLink: fmt.Sprintf("/link/%s/delete", link.ExternalKey),
		})
	}

	return linkViews, nil
}

func (s *State) GetFilesView(linkID string, renderDownloadLinks bool) ([]templates.FileView, error) {
	files, err := s.fs.ListLinkFiles(linkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all files for link %s: %w", linkID, err)
	}

	fileViews := make([]templates.FileView, 0, len(files))
	for _, file := range files {
		downloadLink := ""
		if renderDownloadLinks {
			downloadLink = fmt.Sprintf("/link/%s/file/%s", linkID, file.Name)
		}

		fileViews = append(fileViews, templates.FileView{
			Name:         file.Name,
			UploadedAt:   file.ModTime.Format(time.Stamp),
			Size:         bytesToHuman(file.Size),
			DownloadLink: downloadLink,
			DeleteLink:   fmt.Sprintf("/link/%s/file/%s/delete", linkID, file.Name),
		})
	}

	return fileViews, nil
}

// TODO: don't use fractional decimals for bytes
func bytesToHuman(bytes int64) string {
	b := float64(bytes)
	sizes := []string{"bytes", "KiB", "MiB", "GiB", "TiB", "PiB"}
	i := 0
	for b >= 1024 && i < len(sizes) {
		i++
		b /= 1024
	}
	return fmt.Sprintf("%.2f %s", b, sizes[i])
}
