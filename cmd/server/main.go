package main

import (
	"cmp"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/foxpy/send-me-the-data/cmd/server/handler/admin"
	"github.com/foxpy/send-me-the-data/cmd/server/handler/user"
	"github.com/foxpy/send-me-the-data/cmd/server/idb"
	"github.com/foxpy/send-me-the-data/cmd/server/idb/postgres"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs/vfs"
)

func main() {
	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		slog.Error("required environment variable POSTGRES_URL is not defined")
		os.Exit(1)
	}

	prefix := os.Getenv("PREFIX")
	if prefix == "" {
		slog.Error("required environment variable PREFIX is not defined")
		os.Exit(1)
	}

	userListenAddress := cmp.Or(os.Getenv("USER_LISTEN_ADDRESS"), ":6969")
	adminListenAddress := cmp.Or(os.Getenv("ADMIN_LISTEN_ADDRESS"), ":6767")

	db, err := postgres.NewPostgres(postgresURL)
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}

	fs, err := vfs.NewVFS(prefix)
	if err != nil {
		slog.Error("failed to initalize filesystem", "error", err)
		os.Exit(1)
	}

	err = cleanup(db, fs)
	if err != nil {
		slog.Error("failed to cleanup file journal", "error", err)
		os.Exit(1)
	}

	go func() {
		m := admin.NewAdminServer(db, fs)
		slog.Info("Starting admin HTTP server", "address", adminListenAddress)
		err := http.ListenAndServe(adminListenAddress, m)
		slog.Error("admin ListenAndServe failed", "error", err)
	}()
	go func() {
		m := user.NewUserServer(db, fs)
		slog.Info("Starting user HTTP server", "address", userListenAddress)
		err := http.ListenAndServe(userListenAddress, m)
		slog.Error("user ListenAndServe failed", "error", err)
	}()

	// TODO: graceful shutdown?
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}

func cleanup(db idb.Database, fs ifs.Filesystem) error {
	for {
		entry, err := db.GetFileJournalEntry()
		if errors.Is(err, sql.ErrNoRows) {
			break
		} else if err != nil {
			return fmt.Errorf("failed to obtain a file journal entry: %w", err)
		}

		err = fs.RemoveLinkFile(entry.LinkExternalKey, entry.FileName)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf(
				"failed to delete file %s from link %s referenced by the file journal: %w",
				entry.FileName, entry.LinkExternalKey, err,
			)
		}

		err = db.DeleteFileJournalEntry(entry)
		if err != nil {
			return fmt.Errorf("failed to delete a file journal entry: %w", err)
		}
	}

	return nil
}
