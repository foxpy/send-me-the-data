package postgres

import "github.com/foxpy/send-me-the-data/src/idb"

func (d *Postgres) GetFileJournalEntry() (*idb.FileJournalEntry, error) {
	var entry idb.FileJournalEntry
	err := d.db.QueryRow(`
		SELECT l.public_id, fj.name
		FROM smtd.file_journal fj
			INNER JOIN smtd.links l USING (link_id)
		LIMIT 1
	`).Scan(&entry.LinkPublicID, &entry.FileName)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (d *Postgres) DeleteFileJournalEntry(entry *idb.FileJournalEntry) error {
	_, err := d.db.Exec(`
		DELETE FROM smtd.file_journal
		WHERE name = $1
		  AND link_id = (SELECT link_id FROM smtd.links WHERE public_id = $2)
	`, entry.FileName, entry.LinkPublicID)
	return err
}

func (d *Postgres) CreateFileJournalEntry(entry *idb.FileJournalEntry) error {
	_, err := d.db.Exec(`
		INSERT INTO smtd.file_journal (link_id, name)
		VALUES ((SELECT link_id FROM smtd.links WHERE public_id = $1), $2)
	`, entry.LinkPublicID, entry.FileName)
	return err
}
