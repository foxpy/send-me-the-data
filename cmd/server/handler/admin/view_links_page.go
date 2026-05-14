package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/foxpy/send-me-the-data/cmd/server/flash"
	"github.com/foxpy/send-me-the-data/cmd/server/template"
	"github.com/foxpy/send-me-the-data/cmd/server/view"
)

func (s *AdminServer) viewLinksPage(w http.ResponseWriter, r *http.Request) error {
	offset, err := strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		offset = 0
	}

	totalLinks, err := s.db.TotalLinks()
	if err != nil {
		return fmt.Errorf("failed to count total links: %w", err)
	}

	links, err := s.db.ListLinks(uint(offset), 100)
	if err != nil {
		return fmt.Errorf("failed to query all links from database: %w", err)
	}

	var params template.Params[template.AdminViewLinksParams]
	params.Title = "Send me the Data"
	params.Data.Links = make([]template.LinkView, 0, len(links))
	for _, link := range links {
		files, err := s.fs.ListLinkFiles(link.ID())
		if err != nil {
			return fmt.Errorf("failed to list files for link %s: %w", link.ID(), err)
		}

		params.Data.Links = append(params.Data.Links, view.Link(link, files))
	}
	params.Data.Pages = view.Pagination(uint(offset), uint(totalLinks), 100, "/")

	flashes := flash.GetFlashes(w, r)
	params.SuccessFlash = flashes[flash.SuccessFlash]
	params.ErrorFlash = flashes[flash.ErrorFlash]

	return template.RenderAdminViewLinks(w, params)
}
