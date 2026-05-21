package main

import (
	"cmp"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/foxpy/send-me-the-data/src/handler/admin"
	"github.com/foxpy/send-me-the-data/src/handler/user"
	"github.com/foxpy/send-me-the-data/src/idb"
	"github.com/foxpy/send-me-the-data/src/idb/postgres"
	"github.com/foxpy/send-me-the-data/src/ifs"
	"github.com/foxpy/send-me-the-data/src/ifs/vfs"
	"github.com/foxpy/send-me-the-data/src/irnd/rand"
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

	defaultPasswordHashThreads := (runtime.NumCPU() + 1) / 2
	passwordHashThreads, err := strconv.ParseUint(
		cmp.Or(
			os.Getenv("PASSWORD_HASH_THREADS"),
			strconv.Itoa(defaultPasswordHashThreads),
		),
		10, 64,
	)
	if err != nil || passwordHashThreads == 0 || passwordHashThreads > uint64(runtime.NumCPU())*4 {
		passwordHashThreads = uint64(defaultPasswordHashThreads)
	}

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

	rnd := rand.NewRandom()

	err = cleanup(db, fs)
	if err != nil {
		slog.Error("failed to cleanup file journal", "error", err)
		os.Exit(1)
	}

	adminServer := &http.Server{
		Addr:    adminListenAddress,
		Handler: admin.NewAdminServer(db, fs, rnd, uint(passwordHashThreads)),
		// TODO: set up ReadTimeout or ReadHeaderTimeout
	}
	go func() {
		slog.Info("starting admin HTTP server", "address", adminServer.Addr)
		err := adminServer.ListenAndServe()
		if err != http.ErrServerClosed {
			slog.Error("admin HTTP server failed", "error", err)
		}
	}()

	userServer := &http.Server{
		Addr:    userListenAddress,
		Handler: user.NewUserServer(db, fs),
		// TODO: set up ReadTimeout or ReadHeaderTimeout
	}
	go func() {
		slog.Info("starting user HTTP server", "address", userServer.Addr)
		err := userServer.ListenAndServe()
		if err != http.ErrServerClosed {
			slog.Error("user HTTP server failed", "error", err)
		}
	}()

	// wait for keyboard interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	wg := sync.WaitGroup{}
	wg.Add(2)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	go func() {
		defer wg.Done()
		err := adminServer.Shutdown(ctx)
		if err != nil {
			slog.Error("failed to stop admin HTTP server", "error", err)
		} else {
			slog.Info("admin HTTP server stopped")
		}
	}()

	go func() {
		defer wg.Done()
		err := userServer.Shutdown(ctx)
		if err != nil {
			slog.Error("failed to stop user HTTP server", "error", err)
		} else {
			slog.Info("user HTTP server stopped")
		}
	}()

	wg.Wait()
}

func cleanup(db idb.Database, fs ifs.Filesystem) error {
	for {
		entry, err := db.GetFileJournalEntry()
		if errors.Is(err, sql.ErrNoRows) {
			break
		} else if err != nil {
			return fmt.Errorf("failed to obtain a file journal entry: %w", err)
		}

		err = fs.RemoveLinkFile(entry.LinkPublicID, entry.FileName)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf(
				"failed to delete file %s from link %s referenced by the file journal: %w",
				entry.FileName, entry.LinkPublicID, err,
			)
		}

		err = db.DeleteFileJournalEntry(entry)
		if err != nil {
			return fmt.Errorf("failed to delete a file journal entry: %w", err)
		}
	}

	return nil
}
