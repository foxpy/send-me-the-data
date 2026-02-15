package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/foxpy/send-me-the-data/cmd/server/database"
	"github.com/foxpy/send-me-the-data/cmd/server/filesystem"

	_ "github.com/lib/pq"
)

type State struct {
	db *database.Database
	fs *filesystem.Filesystem
}

func NewState(postgresURL, prefix string) (*State, error) {
	if prefix == "" {
		return nil, errors.New("filesystem prefix must not be empty")
	}

	db, err := database.NewDatabase(postgresURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	fs, err := filesystem.NewFilesystem(prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize filesystem: %w", err)
	}

	return &State{
		db,
		fs,
	}, nil
}

func (s *State) Cleanup() error {
	for {
		path, err := s.db.GetFileJournalEntry()
		if errors.Is(err, sql.ErrNoRows) {
			break
		} else if err != nil {
			return fmt.Errorf("failed to obtain a file journal entry: %w", err)
		}

		err = s.fs.Remove(path)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete a file referenced by the file journal: %s", path)
		}

		err = s.db.DeleteFileJournalEntry(path)
		if err != nil {
			return fmt.Errorf("failed to delete a file journal entry: %s", path)
		}
	}

	return nil
}
