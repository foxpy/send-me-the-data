package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/foxpy/send-me-the-data/cmd/server/idb"
)

type link struct {
	name, externalKey               string
	createdAt                       time.Time
	userDownloadable, uploadEnabled bool
	maxFileSize                     uint64
}

func (l *link) Name() string {
	return l.name
}

func (l *link) ExternalKey() string {
	return l.externalKey
}

func (l *link) CreatedAt() time.Time {
	return l.createdAt
}

func (l *link) UserDownloadable() bool {
	return l.userDownloadable
}

func (l *link) UploadEnabled() bool {
	return l.uploadEnabled
}

func (l *link) MaxFileSize() uint64 {
	return l.maxFileSize
}

type linkRLock struct {
	link
	tx *sql.Tx
}

func (l *linkRLock) Release() error {
	return l.tx.Rollback()
}

type linkWLock struct {
	link
	tx *sql.Tx
}

func (l *linkWLock) Update(name string, userDownloadable, uploadEnabled bool, maxFileSize uint64) error {
	_, err := l.tx.Exec(`
			UPDATE smtd.links SET
				name=$2,
				user_downloadable=$3,
				upload_enabled=$4,
				max_file_size=$5
			WHERE external_key = $1
		`,
		l.externalKey,
		name,
		userDownloadable,
		uploadEnabled,
		maxFileSize,
	)

	if err != nil {
		return fmt.Errorf("failed to update link %s: %w", l.externalKey, err)
	}

	return nil
}

func (l *linkWLock) Delete() error {
	_, err := l.tx.Exec("DELETE FROM smtd.links WHERE external_key = $1", l.externalKey)
	if err != nil {
		return fmt.Errorf("failed to delete link %s: %w", l.externalKey, err)
	}

	return nil
}

func (l *linkWLock) Commit() error {
	return l.tx.Commit()
}

func (l *linkWLock) Rollback() error {
	return l.tx.Rollback()
}

func (d *Postgres) AllLinks() ([]idb.Link, error) {
	rows, err := d.db.Query(`
		SELECT
			name, external_key, created_at, user_downloadable, upload_enabled, max_file_size
		FROM smtd.links
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query links from the database: %w", err)
	}
	defer rows.Close()

	links := make([]idb.Link, 0)
	for rows.Next() {
		var l link
		err = rows.Scan(
			&l.name, &l.externalKey, &l.createdAt, &l.userDownloadable, &l.uploadEnabled, &l.maxFileSize,
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

func (d *Postgres) CreateLink(name, externalKey string, userDownloadable, uploadEnabled bool, maxFileSize uint64) error {
	_, err := d.db.Exec(`
		INSERT INTO smtd.links
			(name, external_key, user_downloadable, upload_enabled, max_file_size)
		VALUES
			($1, $2, $3, $4, $5)`,
		name,
		externalKey,
		userDownloadable,
		uploadEnabled,
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

	var l link
	l.externalKey = externalKey

	err = tx.QueryRow(`
		SELECT
			name, user_downloadable, upload_enabled, created_at, max_file_size
		FROM smtd.links
		WHERE external_key = $1
		FOR SHARE`,
		externalKey,
	).Scan(&l.name, &l.userDownloadable, &l.uploadEnabled, &l.createdAt, &l.maxFileSize)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("failed to acquire read lock on link %s: %w", externalKey, err)
	}

	return &linkRLock{l, tx}, nil
}

func (d *Postgres) AcquireLinkWLock(externalKey string) (idb.LinkWLock, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var l link
	l.externalKey = externalKey
	err = tx.QueryRow(`
		SELECT
			name, user_downloadable, upload_enabled, created_at, max_file_size
		FROM smtd.links
		WHERE external_key = $1
		FOR UPDATE`,
		externalKey,
	).Scan(&l.name, &l.userDownloadable, &l.uploadEnabled, &l.createdAt, &l.maxFileSize)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("failed to acquire write lock on link %s: %w", externalKey, err)
	}

	return &linkWLock{l, tx}, nil
}
