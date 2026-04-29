package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/foxpy/send-me-the-data/cmd/server/idb"
)

type link struct {
	name, id                        string
	createdAt                       time.Time
	userDownloadable, uploadEnabled bool
	maxFileSize                     uint64
}

func (l *link) Name() string {
	return l.name
}

func (l *link) ID() string {
	return l.id
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
			WHERE public_id = $1
		`,
		l.id,
		name,
		userDownloadable,
		uploadEnabled,
		maxFileSize,
	)

	if err != nil {
		return fmt.Errorf("failed to update link %s: %w", l.id, err)
	}

	return nil
}

func (l *linkWLock) Delete() error {
	_, err := l.tx.Exec("DELETE FROM smtd.links WHERE public_id = $1", l.id)
	if err != nil {
		return fmt.Errorf("failed to delete link %s: %w", l.id, err)
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
			name, public_id, created_at, user_downloadable, upload_enabled, max_file_size
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

func (d *Postgres) AcquireLinkRLock(id string) (idb.LinkRLock, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var l link
	l.id = id

	err = tx.QueryRow(`
		SELECT
			name, user_downloadable, upload_enabled, created_at, max_file_size
		FROM smtd.links
		WHERE public_id = $1
		FOR SHARE`,
		id,
	).Scan(&l.name, &l.userDownloadable, &l.uploadEnabled, &l.createdAt, &l.maxFileSize)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("failed to acquire read lock on link %s: %w", id, err)
	}

	return &linkRLock{l, tx}, nil
}

func (d *Postgres) AcquireLinkWLock(id string) (idb.LinkWLock, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var l link
	l.id = id
	err = tx.QueryRow(`
		SELECT
			name, user_downloadable, upload_enabled, created_at, max_file_size
		FROM smtd.links
		WHERE public_id = $1
		FOR UPDATE`,
		id,
	).Scan(&l.name, &l.userDownloadable, &l.uploadEnabled, &l.createdAt, &l.maxFileSize)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("failed to acquire write lock on link %s: %w", id, err)
	}

	return &linkWLock{l, tx}, nil
}
