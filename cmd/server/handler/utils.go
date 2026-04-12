package handler

import (
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/foxpy/send-me-the-data/cmd/server/template"
)

// TODO: implement logging middleware for user and admin servers

//go:embed static/*
var Static embed.FS

type Handler func(http.ResponseWriter, *http.Request) error

func RespondError(w http.ResponseWriter, code int) error {
	w.WriteHeader(code)
	err := template.RenderError(w, code)
	if err != nil {
		slog.Error("failed to write error response", "code", code, "error", err)
	}

	return nil
}

func HandleWith500OnError(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			slog.Error("handler failed", "error", err)
			_ = RespondError(w, http.StatusInternalServerError)
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
