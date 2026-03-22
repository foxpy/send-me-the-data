package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/foxpy/send-me-the-data/cmd/server/handler"
)

func (s *AdminServer) editLink(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	lock, err := s.db.AcquireLinkWLock(id)
	if errors.Is(err, sql.ErrNoRows) {
		return handler.Respond404(w)
	} else if err != nil {
		return fmt.Errorf("failed to acquire write lock for link %s: %w", id, err)
	}

	defer lock.Release()

	name := r.FormValue("name")
	userDownloadable := false
	if r.FormValue("user_downloadable") == "on" {
		userDownloadable = true
	}

	// TODO: check that name is at least not of length 0
	err = lock.UpdateLink(name, userDownloadable)
	if err != nil {
		return fmt.Errorf("failed to update link %s: %w", id, err)
	}

	// TODO: figure out how to handle cookies more nicely
	http.SetCookie(w, &http.Cookie{
		Name:   "success_flash",
		Path:   "/",
		MaxAge: 60,
	})
	http.Redirect(w, r, fmt.Sprintf("/link/%s", id), http.StatusSeeOther)
	return nil
}
