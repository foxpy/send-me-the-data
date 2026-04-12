package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/foxpy/send-me-the-data/cmd/server/flash"
	"github.com/foxpy/send-me-the-data/cmd/server/handler"
)

func (s *AdminServer) editLink(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	maxFileSize, err := strconv.ParseUint(r.FormValue("max_file_size"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	lock, err := s.db.AcquireLinkWLock(id)
	if errors.Is(err, sql.ErrNoRows) {
		return handler.RespondError(w, http.StatusNotFound)
	} else if err != nil {
		return fmt.Errorf("failed to acquire write lock for link %s: %w", id, err)
	}

	defer lock.Release()

	name := r.FormValue("name")
	userDownloadable := false
	if r.FormValue("user_downloadable") == "on" {
		userDownloadable = true
	}

	err = lock.UpdateLink(name, userDownloadable, maxFileSize)
	if err != nil {
		return fmt.Errorf("failed to update link %s: %w", id, err)
	}

	flash.AddFlash(w, flash.SuccessFlash, "Link updated successfully")
	http.Redirect(w, r, fmt.Sprintf("/link/%s", id), http.StatusSeeOther)
	return nil
}
