package admin

import (
	"errors"
	"fmt"
	"net/http"

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
		http.SetCookie(w, &http.Cookie{
			Name:   "error_flash",
			MaxAge: 60,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to create link: %w", err)
	}

	// TODO: set flash text in cookie value EVERYWHERE
	http.SetCookie(w, &http.Cookie{
		Name:   "success_flash",
		MaxAge: 60,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
