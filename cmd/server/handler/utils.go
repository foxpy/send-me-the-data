package handler

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/foxpy/send-me-the-data/cmd/server/idb"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs"
	"github.com/foxpy/send-me-the-data/cmd/server/template"
	"github.com/foxpy/send-me-the-data/cmd/server/view"
)

//go:embed static/*
var Static embed.FS

type Handler func(http.ResponseWriter, *http.Request) error

func Respond404(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNotFound)
	return template.Render404(w)
}

func Respond500(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	err := template.Render500(w)
	if err != nil {
		slog.Error("failed to write 500 response", "error", err)
	}
}

func HandleWith500OnError(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			slog.Error("handler failed", "error", err)
			Respond500(w)
		}
	}
}

func SanitizeFileName(fileName string) (string, error) {
	// TODO: limit file name length somehow:
	//  - decide whether file name limit should be hardcoded or configurable
	//  - maybe I can validate file name length limit from HTML??????????
	//  - decide whether I should return an error or strip file name
	if strings.ContainsAny(fileName, "/\\") {
		return "", fmt.Errorf("Forbidden characters found in file name %s", fileName)
	}

	return fileName, nil
}

func PrepareFilesView(db idb.Database, fs ifs.Filesystem, id string, forAdmin bool) (string, []template.FileView, error) {
	lock, err := db.AcquireLinkRLock(id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil, err
	} else if err != nil {
		return "", nil, fmt.Errorf("failed to acquire read lock on link %s: %w", id, err)
	}

	defer func() {
		_ = lock.Close()
	}()

	renderDownloadLinks := forAdmin || lock.UserDownloadable()
	files, err := view.Files(fs, id, renderDownloadLinks)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get files view for link %s: %w", id, err)
	}

	return lock.Name(), files, nil
}
