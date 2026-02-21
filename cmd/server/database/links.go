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

type LinkRLock struct {
	tx *sql.Tx
}

func (l LinkRLock) Close() error {
	return l.tx.Rollback()
}

type LinkWLock struct {
	externalKey string
	tx          *sql.Tx
}

func (l LinkWLock) Close() error {
	return l.tx.Commit()
}

func (l LinkWLock) DeleteLink() error {
	_, err := l.tx.Exec("DELETE FROM smtd.links WHERE external_key = $1", l.externalKey)
	if err != nil {
		return fmt.Errorf("failed to delete link %s: %w", l.externalKey, err)
	}

	return nil
}

// TODO: maybe I should just take RLock everywhere I use this method
// TODO: delete this method?
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

func (d *Database) AcquireLinkRLock(externalKey string) (*LinkRLock, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	var n int
	err = tx.QueryRow("SELECT 1 FROM smtd.links WHERE external_key = $1 FOR SHARE", externalKey).Scan(&n)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("failed to acquire read lock on link %s: %w", externalKey, err)
	}

	return &LinkRLock{
		tx,
	}, nil
}

func (d *Database) AcquireLinkWLock(externalKey string) (*LinkWLock, error) {
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

	return &LinkWLock{
		externalKey,
		tx,
	}, nil
}
