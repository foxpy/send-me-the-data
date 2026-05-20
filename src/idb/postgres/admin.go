package postgres

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
