package handler

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/foxpy/send-me-the-data/cmd/server/template"
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
	// we accept any file name lengths because names longer than 255 bytes will be rejected by OS anyway

	if strings.ContainsAny(fileName, "/\\") {
		return "", fmt.Errorf("Forbidden characters found in file name %s", fileName)
	}

	return fileName, nil
}
