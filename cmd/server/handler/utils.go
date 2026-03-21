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
	// TODO: limit file name length somehow:
	//  - decide whether file name limit should be hardcoded or configurable
	//  - maybe I can validate file name length limit from HTML??????????
	//  - decide whether I should return an error or strip file name
	if strings.ContainsAny(fileName, "/\\") {
		return "", fmt.Errorf("Forbidden characters found in file name %s", fileName)
	}

	return fileName, nil
}
