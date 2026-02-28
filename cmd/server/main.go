package main

import (
	"cmp"
	"embed"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

//go:embed static/*
var static embed.FS

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

	state, err := NewState(postgresURL, prefix)
	if err != nil {
		slog.Error("failed to initialize the app", "error", err)
		os.Exit(1)
	}

	err = state.Cleanup()
	if err != nil {
		slog.Error("failed to cleanup file journal", "error", err)
		os.Exit(1)
	}

	go adminServer(state, adminListenAddress)
	go userServer(state, userListenAddress)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}

func userServer(state *State, listenAddress string) {
	m := http.NewServeMux()
	m.HandleFunc("GET /u/{id}", handleWith500OnError(state.handleUserViewLinkPage))
	m.HandleFunc("POST /u/{id}", handleWith500OnError(state.handleUserUpload))
	m.HandleFunc("GET /link/{id}/file/{name}", handleWith500OnError(state.handleUserDownloadFile))
	m.Handle("GET /static/", http.FileServerFS(static))

	slog.Info("Starting user HTTP server", "address", listenAddress)
	err := http.ListenAndServe(listenAddress, m)
	slog.Error("user ListenAndServe failed", "error", err)
}

func adminServer(state *State, listenAddress string) {
	m := http.NewServeMux()
	m.HandleFunc("GET /{$}", handleWith500OnError(state.handleAdminViewLinksPage))
	m.HandleFunc("GET /link/{id}", handleWith500OnError(state.handleAdminViewLinkPage))
	m.HandleFunc("GET /link/{id}/file/{name}", handleWith500OnError(state.handleAdminDownloadFile))
	m.HandleFunc("POST /link/{id}/file/{name}/delete", handleWith500OnError(state.handleAdminDeleteFile))
	m.HandleFunc("POST /link/{id}/delete", handleWith500OnError(state.handleAdminDeleteLink))
	m.HandleFunc("POST /link", handleWith500OnError(state.handleAdminCreateLink))
	m.Handle("GET /static/", http.FileServerFS(static))

	slog.Info("Starting admin HTTP server", "address", listenAddress)
	err := http.ListenAndServe(listenAddress, m)
	slog.Error("admin ListenAndServe failed", "error", err)
}
