package main

import (
	_ "embed"
	"net/http"
)

//go:embed templates/500.gohtml
var internalServerErrorTemplate string

func respond500(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(internalServerErrorTemplate))
}
