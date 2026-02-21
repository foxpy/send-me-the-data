package main

import (
	"fmt"
	"net/http"
)

func (s *State) handleDeleteLink(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	ok, err := s.db.DoesLinkExist(id)
	if err != nil {
		return fmt.Errorf("failed to check if link is published: %w", err)
	}

	if !ok {
		return respond404(w)
	}

	// FIXME: since we do not attempt to terminate all active uploads,
	//        acquiring this lock might take possibly many hours of time.
	lock, err := s.db.AcquireLinkWLock(id)
	if err != nil {
		return fmt.Errorf("failed to acquire write lock for link %s: %w", id, err)
	}

	defer func() {
		_ = lock.Close()
	}()

	err = s.fs.RemoveLinkFiles(id)
	if err != nil {
		return err
	}

	err = lock.DeleteLink()
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
