package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

func (s *State) handleAdminDeleteFile(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	lock, err := s.db.AcquireLinkRLock(id)
	if errors.Is(err, sql.ErrNoRows) {
		return respond404(w)
	} else if err != nil {
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
