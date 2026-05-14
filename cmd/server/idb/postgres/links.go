package postgres

import (
	"fmt"

	"github.com/foxpy/send-me-the-data/cmd/server/idb"
)

func (d *Postgres) TotalLinks() (uint64, error) {
	row := d.db.QueryRow(`
		SELECT count(link_id) FROM smtd.links
	`)

	var n uint64
	err := row.Scan(&n)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (d *Postgres) ListLinks(offset, limit uint) ([]idb.Link, error) {
	rows, err := d.db.Query(`
		SELECT
			name, public_id, created_at, user_downloadable, upload_enabled, max_file_size
		FROM smtd.links
		ORDER BY created_at DESC
		LIMIT $1
		OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query links from the database: %w", err)
	}
	defer rows.Close()

	links := make([]idb.Link, 0)
	for rows.Next() {
		var l link
		err = rows.Scan(
			&l.name, &l.id, &l.createdAt, &l.userDownloadable, &l.uploadEnabled, &l.maxFileSize,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan link from the database: %w", err)
		}

		links = append(links, &l)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over links from the database: %w", err)
	}

	return links, nil
}

func (d *Postgres) CreateLink(name, id string, userDownloadable, uploadEnabled bool, maxFileSize uint64) error {
	_, err := d.db.Exec(`
		INSERT INTO smtd.links
			(name, public_id, user_downloadable, upload_enabled, max_file_size)
		VALUES
			($1, $2, $3, $4, $5)`,
		name,
		id,
		userDownloadable,
		uploadEnabled,
		maxFileSize,
	)
	if err != nil {
		return fmt.Errorf("failed to create new link %s: %w", id, err)
	}

	return nil
}
