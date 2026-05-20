package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/foxpy/send-me-the-data/src/idb"
)

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

func (l *linkWLock) Rollback() {
	err := l.tx.Rollback()
	if err != nil && !errors.Is(err, sql.ErrTxDone) {
		slog.Error("failed to release link wlock", "link", l.id, "error", err)
	}
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
