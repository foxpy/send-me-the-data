package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/foxpy/send-me-the-data/cmd/server/flash"
	"github.com/foxpy/send-me-the-data/cmd/server/handler"
	"github.com/foxpy/send-me-the-data/cmd/server/template"
	"github.com/foxpy/send-me-the-data/cmd/server/view"
)

func (s *AdminServer) viewLinkPage(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	lock, err := s.db.AcquireLinkRLock(id)
	if errors.Is(err, sql.ErrNoRows) {
		return handler.RespondError(w, http.StatusNotFound)
	} else if err != nil {
		return fmt.Errorf("failed to acquire read lock on link %s: %w", id, err)
	}

	defer lock.Release()

	files, err := view.Files(s.fs, lock)
	if err != nil {
		return fmt.Errorf("failed to get files view for link %s: %w", id, err)
	}

	link, err := view.Link(lock, s.fs)
	if err != nil {
		return fmt.Errorf("failed to get %s link view: %w", id, err)
	}

	lock.Release()

	var params template.Params[template.AdminViewLinkParams]
	params.Title = fmt.Sprintf("Link: %s", lock.Name())
	params.Data.Files = files
	params.Data.Link = *link

	flashes := flash.GetFlashes(w, r)
	params.SuccessFlash = flashes[flash.SuccessFlash]
	params.ErrorFlash = flashes[flash.ErrorFlash]

	return template.RenderAdminViewLink(w, params)
}
