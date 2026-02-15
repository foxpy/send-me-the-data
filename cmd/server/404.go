package main

import (
	"net/http"

	_ "embed"
)

//go:embed templates/404.gohtml
var notFoundTemplate string

func respond404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(notFoundTemplate))
}
