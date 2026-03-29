package admin

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/foxpy/send-me-the-data/cmd/server/flash"
	"github.com/lib/pq"
)

func (s *AdminServer) createLink(w http.ResponseWriter, r *http.Request) error {
	name := r.FormValue("name")
	externalKey := s.db.GenerateRandomExternalKey()
	userDownloadable := false
	if r.FormValue("user_downloadable") == "on" {
		userDownloadable = true
	}

	err := s.db.CreateLink(name, externalKey, userDownloadable)
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
