package postgres

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/foxpy/send-me-the-data/cmd/server/idb"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type Postgres struct {
	db *sql.DB
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

var _ idb.Database = &Postgres{}

func NewPostgres(postgresURL string) (*Postgres, error) {
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

	return &Postgres{
		db,
	}, nil
}
