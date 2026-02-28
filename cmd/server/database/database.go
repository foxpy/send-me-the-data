package database

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Database struct {
	db *sql.DB
}

func NewDatabase(postgresURL string) (*Database, error) {
	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	goose.SetBaseFS(embedMigrations)
	err = goose.SetDialect("postgres")
	if err != nil {
		return nil, fmt.Errorf("failed to set goose dialect: %w", err)
	}

	err = goose.Up(db, "migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to apply database migrations: %w", err)
	}

	return &Database{
		db,
	}, nil
}
