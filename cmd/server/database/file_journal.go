package database

func (d *Database) GetFileJournalEntry() (string, error) {
	var path string
	err := d.db.QueryRow("SELECT path FROM smtd.file_journal LIMIT 1").Scan(&path)
	if err != nil {
		return "", err
	}

	return path, nil
}

func (d *Database) DeleteFileJournalEntry(path string) error {
	_, err := d.db.Exec("DELETE FROM smtd.file_journal WHERE path = $1", path)
	return err
}

func (d *Database) CreateFileJournalEntry(path string) error {
	_, err := d.db.Exec("INSERT INTO smtd.file_journal VALUES ($1)", path)
	return err
}
