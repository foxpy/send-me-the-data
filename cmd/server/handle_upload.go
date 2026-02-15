package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

func (s *State) handleUpload(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	ok, err := doesLinkExist(s.db, id)
	if err != nil {
		return fmt.Errorf("failed to check if link is published: %w", err)
	}

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(notFoundTemplate))
		return nil
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return fmt.Errorf("failed to get file from form data: %w", err)
	}

	err = os.MkdirAll(fmt.Sprintf("%s/%s", s.prefix, id), 0777)
	if err != nil {
		return fmt.Errorf("failed to create a directory for saved files: %w", err)
	}

	filePath := fmt.Sprintf("%s/%s/%s", s.prefix, id, header.Filename)

	_, err = s.db.Exec("INSERT INTO smtd.file_journal VALUES ($1)", filePath)
	if err != nil {
		return fmt.Errorf("failed to create a file journal entry: %w", err)
	}

	defer func() {
		_, err = s.db.Exec("DELETE FROM smtd.file_journal WHERE path = $1", filePath)
		if err != nil {
			slog.Error("failed to remove a file journal entry", "path", filePath, "error", err)
		}
	}()

	dirty := true
	localFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}

	defer func() {
		_ = localFile.Close()

		if dirty {
			err := os.Remove(filePath)
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
