package main

import (
	"fmt"
	"net/http"
)

func (s *State) handleAdminDeleteFile(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	ok, err := s.db.DoesLinkExist(id)
	if err != nil {
		return fmt.Errorf("failed to check if link is published: %w", err)
	}

	if !ok {
		return respond404(w)
	}

	lock, err := s.db.AcquireLinkRLock(id)
	if err != nil {
		return fmt.Errorf("failed to acquire read lock for link %s: %w", id, err)
	}

	defer func() {
		_ = lock.Close()
	}()

	name := r.PathValue("name")
	err = s.fs.RemoveLinkFile(id, name)
	if err != nil {
		return fmt.Errorf("failed to remove file %s from link %s: %w", name, id, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/link/%s", id), http.StatusSeeOther)
	return nil
}
