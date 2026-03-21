package admin

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/foxpy/send-me-the-data/cmd/server/handler"
	"github.com/foxpy/send-me-the-data/cmd/server/template"
)

func (s *AdminServer) viewLinkPage(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	title, files, err := handler.PrepareFilesView(s.db, s.fs, id, true)
	if errors.Is(err, sql.ErrNoRows) {
		return handler.Respond404(w)
	} else if err != nil {
		return err
	}

	var params template.Params[template.AdminViewLinkParams]
	params.Title = title
	params.Data.Files = files

	return template.RenderAdminViewLink(w, params)
}
