package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/foxpy/send-me-the-data/cmd/server/idb"
)

type linkRLock struct {
	name, externalKey string
	userDownloadable  bool
	createdAt         time.Time
	maxFileSize       uint64
	tx                *sql.Tx
}

func (l linkRLock) UserDownloadable() bool {
	return l.userDownloadable
}

func (l linkRLock) Name() string {
	return l.name
}

func (l linkRLock) ExternalKey() string {
	return l.externalKey
}

func (l linkRLock) CreatedAt() time.Time {
	return l.createdAt
}

func (l linkRLock) MaxFileSize() uint64 {
	return l.maxFileSize
}

func (l linkRLock) Release() error {
	return l.tx.Rollback()
}

type linkWLock struct {
	externalKey string
	tx          *sql.Tx
}

func (l linkWLock) Release() error {
	return l.tx.Commit()
}

func (l linkWLock) DeleteLink() error {
	_, err := l.tx.Exec("DELETE FROM smtd.links WHERE external_key = $1", l.externalKey)
	if err != nil {
		return fmt.Errorf("failed to delete link %s: %w", l.externalKey, err)
	}

	return nil
}

func (l linkWLock) UpdateLink(name string, userDownloadable bool, maxFileSize uint64) error {
	_, err := l.tx.Exec(
		"UPDATE smtd.links SET name=$2, user_downloadable=$3, max_file_size=$4 WHERE external_key = $1",
		l.externalKey,
		name,
		userDownloadable,
		maxFileSize,
	)
	if err != nil {
		return fmt.Errorf("failed to update link %s: %w", l.externalKey, err)
	}

	return nil
}

func (d *Postgres) AllLinks() ([]idb.Link, error) {
	rows, err := d.db.Query("SELECT name, external_key, created_at, user_downloadable, max_file_size FROM smtd.links")
	if err != nil {
		return nil, fmt.Errorf("failed to query links from the database: %w", err)
	}
	defer rows.Close()

	links := make([]idb.Link, 0)
	for rows.Next() {
		var link idb.Link
		err = rows.Scan(&link.Name, &link.ExternalKey, &link.CreatedAt, &link.UserDownloadable, &link.MaxFileSize)
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

func (d *Postgres) CreateLink(name, externalKey string, userDownloadable bool, maxFileSize uint64) error {
	_, err := d.db.Exec(
		"INSERT INTO smtd.links (name, external_key, user_downloadable, max_file_size) VALUES ($1, $2, $3, $4)",
		name,
		externalKey,
		userDownloadable,
		maxFileSize,
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

	var lock linkRLock
	lock.tx = tx
	lock.externalKey = externalKey

	err = tx.QueryRow(
		"SELECT name, user_downloadable, created_at, max_file_size FROM smtd.links WHERE external_key = $1 FOR SHARE",
		externalKey,
	).Scan(&lock.name, &lock.userDownloadable, &lock.createdAt, &lock.maxFileSize)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("failed to acquire read lock on link %s: %w", externalKey, err)
	}

	return &lock, nil
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
