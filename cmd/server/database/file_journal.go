package database

type FileJournalEntry struct {
	LinkExternalKey string
	FileName        string
}

func (d *Database) GetFileJournalEntry() (*FileJournalEntry, error) {
	var entry FileJournalEntry
	err := d.db.QueryRow(`
		SELECT l.external_key, fj.name
		FROM smtd.file_journal fj
			INNER JOIN smtd.links l USING (link_id)
		LIMIT 1
	`).Scan(&entry.LinkExternalKey, &entry.FileName)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (d *Database) DeleteFileJournalEntry(entry *FileJournalEntry) error {
	_, err := d.db.Exec(`
		DELETE FROM smtd.file_journal
		WHERE name = $1
		  AND link_id = (SELECT link_id FROM smtd.links WHERE external_key = $2)
	`, entry.FileName, entry.LinkExternalKey)
	return err
}

func (d *Database) CreateFileJournalEntry(entry *FileJournalEntry) error {
	_, err := d.db.Exec(`
		INSERT INTO smtd.file_journal (link_id, name)
		VALUES ((SELECT link_id FROM smtd.links WHERE external_key = $1), $2)
	`, entry.LinkExternalKey, entry.FileName)
	return err
}
