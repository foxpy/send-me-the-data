package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/foxpy/send-me-the-data/cmd/server/handler"
	"github.com/foxpy/send-me-the-data/cmd/server/template"
	"github.com/foxpy/send-me-the-data/cmd/server/view"
)

func (s *AdminServer) editLinkPage(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	lock, err := s.db.AcquireLinkRLock(id)
	if errors.Is(err, sql.ErrNoRows) {
		return handler.Respond404(w)
	} else if err != nil {
		return fmt.Errorf("failed to acquire read lock on link %s: %w", id, err)
	}

	defer func() {
		_ = lock.Release()
	}()

	linkView, err := view.Link(lock, s.fs)
	if err != nil {
		return fmt.Errorf("failed to render link view: %w", err)
	}

	var params template.Params[template.AdminEditLinkParams]
	params.Title = fmt.Sprintf("Edit link: %s", lock.Name())
	params.Data.Link = *linkView

	return template.RenderAdminEditLink(w, params)
}
