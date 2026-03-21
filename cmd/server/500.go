package main

import (
	"log/slog"
	"net/http"

	"github.com/foxpy/send-me-the-data/cmd/server/template"
)

func respond500(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	err := template.Render500(w)
	if err != nil {
		slog.Error("failed to write 500 response", "error", err)
	}
}
