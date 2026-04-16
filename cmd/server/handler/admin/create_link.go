package admin

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/foxpy/send-me-the-data/cmd/server/flash"
	"github.com/lib/pq"
)

func (s *AdminServer) createLink(w http.ResponseWriter, r *http.Request) error {
	name := r.FormValue("name")
	maxFileSize, err := strconv.ParseUint(r.FormValue("max_file_size"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

	externalKey := s.db.GenerateRandomExternalKey()
	userDownloadable := false
	if r.FormValue("user_downloadable") == "on" {
		userDownloadable = true
	}
	uploadEnabled := false
	if r.FormValue("upload_enabled") == "on" {
		uploadEnabled = true
	}

	err = s.db.CreateLink(name, externalKey, userDownloadable, uploadEnabled, maxFileSize)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation" {
		// hopefully should never happen because externalKey collision chance is very low
		flash.AddFlash(w, flash.ErrorFlash, "Failed to create link, try again")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to create link: %w", err)
	}

	flash.AddFlash(w, flash.SuccessFlash, "Link created successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
