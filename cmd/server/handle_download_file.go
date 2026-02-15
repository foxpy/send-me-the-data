package main

import (
	"fmt"
	"net/http"
)

func (s *State) handleDownloadFile(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	ok, err := doesLinkExist(s.db, id)
	if err != nil {
		return fmt.Errorf("failed to check if link is published: %w", err)
	}

	if !ok {
		respond404(w)
		return nil
	}

	name := r.PathValue("name")
	path := fmt.Sprintf("%s/%s/%s", s.prefix, id, name)
	http.ServeFile(w, r, path)
	return nil
}
