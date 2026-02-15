package main

import (
	"fmt"
	"net/http"

	_ "embed"
)

//go:embed templates/404.gohtml
var notFoundTemplate string

func respond404(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte(notFoundTemplate))
	if err != nil {
		return fmt.Errorf("failed to write 404 response: %w", err)
	}

	return nil
}
