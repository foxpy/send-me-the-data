package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/foxpy/send-me-the-data/cmd/server/templates"
)

func (s *State) handleAdminViewLinkPage(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	title, files, err := s.prepareFilesView(id, true)
	if errors.Is(err, sql.ErrNoRows) {
		return respond404(w)
	} else if err != nil {
		return err
	}

	var params templates.Params[templates.AdminViewLinkParams]
	params.Title = title
	params.Data.Files = files

	return templates.RenderAdminViewLink(w, params)
}
