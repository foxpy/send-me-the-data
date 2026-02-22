package main

import (
	"fmt"
	"net/http"
)

func (s *State) handleAdminDownloadFile(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	ok, err := s.db.DoesLinkExist(id)
	if err != nil {
		return fmt.Errorf("failed to check if link is published: %w", err)
	}

	if !ok {
		return respond404(w)
	}

	// TODO: figure out how to sanitize file name
	name := r.PathValue("name")

	fs, err := s.fs.FS(id)
	if err != nil {
		return fmt.Errorf("failed to open filesystem to serve file: %w", err)
	}

	http.ServeFileFS(w, r, fs, name)
	return nil
}
