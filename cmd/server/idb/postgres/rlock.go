package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/foxpy/send-me-the-data/cmd/server/idb"
)

type linkRLock struct {
	link
	tx *sql.Tx
}

func (l *linkRLock) Release() {
	err := l.tx.Rollback()
	if err != nil && !errors.Is(err, sql.ErrTxDone) {
		slog.Error("failed to release link rlock", "link", l.id, "error", err)
	}
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
