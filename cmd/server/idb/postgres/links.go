package postgres

import (
	"database/sql"
	"fmt"

	"github.com/foxpy/send-me-the-data/cmd/server/idb"
)

type linkRLock struct {
	name             string
	userDownloadable bool
	tx               *sql.Tx
}

func (l linkRLock) UserDownloadable() bool {
	return l.userDownloadable
}

func (l linkRLock) Name() string {
	return l.name
}

func (l linkRLock) Close() error {
	return l.tx.Rollback()
}

type linkWLock struct {
	externalKey string
	tx          *sql.Tx
}

func (l linkWLock) Close() error {
	return l.tx.Commit()
}

func (l linkWLock) DeleteLink() error {
	_, err := l.tx.Exec("DELETE FROM smtd.links WHERE external_key = $1", l.externalKey)
	if err != nil {
		return fmt.Errorf("failed to delete link %s: %w", l.externalKey, err)
	}

	return nil
}

func (d *Postgres) AllLinks() ([]idb.Link, error) {
	rows, err := d.db.Query("SELECT name, external_key, created_at FROM smtd.links")
	if err != nil {
		return nil, fmt.Errorf("failed to query links from the database: %w", err)
	}
	defer rows.Close()

	links := make([]idb.Link, 0)
	for rows.Next() {
		var link idb.Link
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

func (d *Postgres) CreateLink(name, externalKey string, userDownloadable bool) error {
	_, err := d.db.Exec(
		"INSERT INTO smtd.links (name, external_key, user_downloadable) VALUES ($1, $2, $3)",
		name,
		externalKey,
		userDownloadable,
	)
	if err != nil {
		return fmt.Errorf("failed to create new link %s: %w", externalKey, err)
	}

	return nil
}

func (d *Postgres) AcquireLinkRLock(externalKey string) (idb.LinkRLock, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var userDownloadable bool
	var name string
	err = tx.QueryRow(
		"SELECT name, user_downloadable FROM smtd.links WHERE external_key = $1 FOR SHARE",
		externalKey,
	).Scan(&name, &userDownloadable)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("failed to acquire read lock on link %s: %w", externalKey, err)
	}

	return &linkRLock{
		name,
		userDownloadable,
		tx,
	}, nil
}

func (d *Postgres) AcquireLinkWLock(externalKey string) (idb.LinkWLock, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var n int
	err = tx.QueryRow("SELECT 1 FROM smtd.links WHERE external_key = $1 FOR UPDATE", externalKey).Scan(&n)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("failed to acquire write lock on link %s: %w", externalKey, err)
	}

	return &linkWLock{
		externalKey,
		tx,
	}, nil
}
