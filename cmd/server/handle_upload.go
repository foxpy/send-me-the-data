package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func (s *State) handleUpload(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	ok, err := s.db.DoesLinkExist(id)
	if err != nil {
		return fmt.Errorf("failed to check if link is published: %w", err)
	}

	if !ok {
		return respond404(w)
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return fmt.Errorf("failed to get file from form data: %w", err)
	}

	lock, err := s.db.AcquireLinkRLock(id)
	if err != nil {
		return fmt.Errorf("failed to acquire read lock on link %s: %w", id, err)
	}

	defer func() {
		_ = lock.Close()
	}()

	// TODO: should we somehow validate/sanitize header.Filename?
	filePath := s.fs.GetPath(id, header.Filename)
	err = s.db.CreateFileJournalEntry(filePath)
	if err != nil {
		return fmt.Errorf("failed to create a file journal entry: %w", err)
	}

	defer func() {
		err = s.db.DeleteFileJournalEntry(filePath)
		if err != nil {
			slog.Error("failed to remove a file journal entry", "path", filePath, "error", err)
		}
	}()

	dirty := true
	localFile, err := s.fs.CreateNewFile(id, header.Filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}

	defer func() {
		_ = localFile.Close()

		if dirty {
			err := s.fs.Remove(filePath)
			if err != nil {
				slog.Error("failed to delete a file", "path", filePath, "error", err)
			}
		}
	}()

	_, err = io.Copy(localFile, file)
	if err != nil {
		return fmt.Errorf("failed to save file %s to storage: %w", filePath, err)
	}

	dirty = false

	http.Redirect(w, r, fmt.Sprintf("/u/%s", id), http.StatusSeeOther)
	return nil
}
