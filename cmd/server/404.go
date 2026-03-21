package main

import (
	"net/http"

	"github.com/foxpy/send-me-the-data/cmd/server/template"
)

func respond404(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNotFound)
	return template.Render404(w)
}
