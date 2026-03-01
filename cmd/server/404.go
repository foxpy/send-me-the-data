package main

import (
	"net/http"

	"github.com/foxpy/send-me-the-data/cmd/server/templates"
)

func respond404(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNotFound)
	return templates.Render404(w)
}
