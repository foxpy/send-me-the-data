package admin

import (
	"fmt"
	"net/http"

	"github.com/foxpy/send-me-the-data/cmd/server/flash"
	"github.com/foxpy/send-me-the-data/cmd/server/template"
	"github.com/foxpy/send-me-the-data/cmd/server/view"
)

func (s *AdminServer) viewLinksPage(w http.ResponseWriter, r *http.Request) error {
	links, err := view.Links(s.db, s.fs)
	if err != nil {
		return fmt.Errorf("failed to get links view: %w", err)
	}

	var params template.Params[template.AdminViewLinksParams]
	params.Title = "Send me the Data"
	params.Data.Links = links

	flashes := flash.GetFlashes(w, r)
	params.SuccessFlash = flashes[flash.SuccessFlash]
	params.ErrorFlash = flashes[flash.ErrorFlash]

	return template.RenderAdminViewLinks(w, params)
}
