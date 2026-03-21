package user

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/foxpy/send-me-the-data/cmd/server/handler"
	"github.com/foxpy/send-me-the-data/cmd/server/idb"
)

func (s *UserServer) upload(w http.ResponseWriter, r *http.Request) error {
	// TODO: idea: can I send the file within a body?
	// FIXME: call to FormFile() actually creates a temporary file and blocks until everything is downloaded
	file, header, err := r.FormFile("file")
	if err != nil {
		return fmt.Errorf("failed to get file from form data: %w", err)
	}

	defer file.Close()
	fileName, err := handler.SanitizeFileName(header.Filename)
	if err != nil {
		return err
	}

	id := r.PathValue("id")
	lock, err := s.db.AcquireLinkRLock(id)
	if errors.Is(err, sql.ErrNoRows) {
		return handler.Respond404(w)
	} else if err != nil {
		return fmt.Errorf("failed to acquire read lock on link %s: %w", id, err)
	}

	defer func() {
		_ = lock.Close()
	}()

	fileJournalEntry := &idb.FileJournalEntry{
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
