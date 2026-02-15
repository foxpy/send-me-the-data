package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type State struct {
	db     *sql.DB
	prefix string
}

func NewState(postgresURL string) (*State, error) {
	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	prefix := os.Getenv("PREFIX")
	if prefix == "" {
		return nil, fmt.Errorf("PREFIX must be provided: %w", err)
	}

	return &State{
		db,
		prefix,
	}, nil
}

func (s *State) Cleanup() error {
	for {
		var path string
		err := s.db.QueryRow("SELECT path FROM smtd.file_journal LIMIT 1").Scan(&path)
		if errors.Is(err, sql.ErrNoRows) {
			break
		} else if err != nil {
			return fmt.Errorf("failed to obtain a file journal entriy: %w", err)
		}

		err = os.Remove(path)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete a file from a file journal: %s", path)
		}

		_, err = s.db.Exec("DELETE FROM smtd.file_journal WHERE path = $1", path)
		if err != nil {
			return fmt.Errorf("failed to delete a file journal entry: %s", path)
		}
	}

	return nil
}
