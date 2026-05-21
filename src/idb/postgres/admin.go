package postgres

import (
	"fmt"
)

func (d *Postgres) GetAdminPasswordHash(username string) ([]byte, error) {
	var hash []byte
	err := d.db.QueryRow(`
		SELECT password_hash FROM smtd.admins WHERE username = $1
	`, username).Scan(&hash)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func (d *Postgres) CreateAdmin(username string, passwordHash []byte) error {
	_, err := d.db.Exec(`
		INSERT INTO smtd.admins (username, password_hash) VALUES ($1, $2)
	`, username, passwordHash)
	if err != nil {
		return fmt.Errorf("failed to create admin: %w", err)
	}

	return nil
}

func (d *Postgres) DeleteAdmin(username string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var adminID int
	err = tx.QueryRow(`
		SELECT admin_id FROM smtd.admins WHERE username = $1 FOR UPDATE
	`, username).Scan(&adminID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM smtd.session_tokens WHERE admin_id = $1`, adminID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM smtd.admins WHERE admin_id = $1`, adminID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
