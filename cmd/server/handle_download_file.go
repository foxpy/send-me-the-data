package main

import (
	"fmt"
	"net/http"
)

func (s *State) handleDownloadFile(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	ok, err := s.db.DoesLinkExist(id)
	if err != nil {
		return fmt.Errorf("failed to check if link is published: %w", err)
	}

	if !ok {
		return respond404(w)
	}

	name := r.PathValue("name")
	// TODO: figure out how to sanitize file name
	http.ServeFileFS(w, r, s.fs.FS(), name)
	return nil
}
