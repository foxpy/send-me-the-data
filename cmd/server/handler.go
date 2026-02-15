package main

import (
	"log/slog"
	"net/http"
)

type Handler func(http.ResponseWriter, *http.Request) error

func handleWith500OnError(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			slog.Error("handler failed", "error", err)
			respond500(w)
		}
	}
}
