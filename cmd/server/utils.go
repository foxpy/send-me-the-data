package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"
)

type File struct {
	Name       string
	UploadedAt string
	Size       string
	// FIXME: learn what fucking MVC is
	size         int64
	DownloadLink string
}

type Link struct {
	Name          string
	CreatedAt     string
	TotalFiles    int
	TotalSize     string
	ViewLinkURL   string
	DeleteLinkURL string
}

func doesLinkExist(db *sql.DB, id string) (bool, error) {
	var n int
	err := db.QueryRow("SELECT 1 FROM smtd.links WHERE external_key = $1", id).Scan(&n)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func listLinks(prefix string, db *sql.DB) ([]Link, error) {
	// FIXME: do not read all links from database, use pagination
	rows, err := db.Query("SELECT name, external_key, created_at FROM smtd.links")
	if err != nil {
		return nil, fmt.Errorf("failed to query links from the database: %w", err)
	}
	defer rows.Close()

	links := make([]Link, 0)
	var name, externalKey string
	var createdAt time.Time
	for rows.Next() {
		err = rows.Scan(&name, &externalKey, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan link from the database: %w", err)
		}

		files, err := listFiles(prefix, externalKey)
		if err != nil {
			return nil, fmt.Errorf("failed to list files for link %s: %w", externalKey, err)
		}

		var totalSize int64
		for _, file := range files {
			totalSize += file.size
		}

		links = append(links, Link{
			Name:          name,
			CreatedAt:     createdAt.Format(time.Stamp),
			TotalFiles:    len(files),
			TotalSize:     bytesToHuman(totalSize),
			ViewLinkURL:   fmt.Sprintf("/u/%s", externalKey),
			DeleteLinkURL: fmt.Sprintf("/delete/%s", externalKey),
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over links from the database: %w", err)
	}

	return links, nil
}

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

func listFiles(prefix, id string) ([]File, error) {
	// FIXME: do not read all dir entries into memory at once
	// solution: pagination
	entries, err := os.ReadDir(fmt.Sprintf("%s/%s", prefix, id))
	if err != nil {
		return nil, fmt.Errorf("failed to traverse directory: %w", err)
	}

	files := make([]File, 0, len(entries))

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("failed to get file info: %w", err)
		}

		files = append(files, File{
			Name:         entry.Name(),
			UploadedAt:   info.ModTime().Format(time.Stamp),
			Size:         bytesToHuman(info.Size()),
			size:         info.Size(),
			DownloadLink: fmt.Sprintf("/f/%s/%s", id, entry.Name()),
		})
	}

	return files, nil
}
