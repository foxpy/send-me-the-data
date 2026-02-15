package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Link struct {
	Name, ExternalKey string
	CreatedAt         time.Time
}

func (d *Database) DoesLinkExist(externalKey string) (bool, error) {
	var n int
	err := d.db.QueryRow("SELECT 1 FROM smtd.links WHERE external_key = $1", externalKey).Scan(&n)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (d *Database) AllLinks() ([]Link, error) {
	// FIXME: do not read all links from database, use pagination
	rows, err := d.db.Query("SELECT name, external_key, created_at FROM smtd.links")
	if err != nil {
		return nil, fmt.Errorf("failed to query links from the database: %w", err)
	}
	defer rows.Close()

	links := make([]Link, 0)
	for rows.Next() {
		var link Link
		err = rows.Scan(&link.Name, &link.ExternalKey, &link.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan link from the database: %w", err)
		}

		links = append(links, link)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over links from the database: %w", err)
	}

	return links, nil
}
