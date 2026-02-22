package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/foxpy/send-me-the-data/cmd/server/database"
)

func (s *State) handleUserUpload(w http.ResponseWriter, r *http.Request) error {
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

	fileName, err := sanitizeFileName(header.Filename)
	if err != nil {
		return err
	}

	lock, err := s.db.AcquireLinkRLock(id)
	if err != nil {
		return fmt.Errorf("failed to acquire read lock on link %s: %w", id, err)
	}

	defer func() {
		_ = lock.Close()
	}()

	fileJournalEntry := &database.FileJournalEntry{
		LinkExternalKey: id,
		FileName:        fileName,
	}
	err = s.db.CreateFileJournalEntry(fileJournalEntry)
	if err != nil {
		return fmt.Errorf("failed to create a file journal entry: %w", err)
	}

	defer func() {
		err = s.db.DeleteFileJournalEntry(fileJournalEntry)
		if err != nil {
			slog.Error("failed to remove a file journal entry", "entry", fileJournalEntry, "error", err)
		}
	}()

	dirty := true
	localFile, err := s.fs.CreateNewFile(id, fileName)
	if err != nil {
		return fmt.Errorf("failed to create file %s for link %s: %w", fileName, id, err)
	}

	defer func() {
		_ = localFile.Close()

		if dirty {
			err := s.fs.RemoveLinkFile(id, fileName)
			if err != nil {
				slog.Error("failed to delete a file", "name", fileName, "link_id", id, "error", err)
			}
		}
	}()

	_, err = io.Copy(localFile, file)
	if err != nil {
		return fmt.Errorf("failed to save file %s from link %s to storage: %w", fileName, id, err)
	}

	dirty = false

	http.SetCookie(w, &http.Cookie{
		Name:   "success_flash",
		MaxAge: 60,
	})
	http.Redirect(w, r, fmt.Sprintf("/u/%s", id), http.StatusSeeOther)
	return nil
}

func sanitizeFileName(fileName string) (string, error) {
	if strings.ContainsAny(fileName, "/\\") {
		return "", fmt.Errorf("Forbidden characters found in file name %s", fileName)
	}

	return fileName, nil
}
