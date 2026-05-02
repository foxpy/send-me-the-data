package user

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/foxpy/send-me-the-data/cmd/server/flash"
	"github.com/foxpy/send-me-the-data/cmd/server/handler"
	"github.com/foxpy/send-me-the-data/cmd/server/template"
	"github.com/foxpy/send-me-the-data/cmd/server/view"
)

func (s *UserServer) viewLinkPage(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	lock, err := s.db.AcquireLinkRLock(id)
	if errors.Is(err, sql.ErrNoRows) {
		return handler.RespondError(w, http.StatusNotFound)
	} else if err != nil {
		return fmt.Errorf("failed to acquire read lock on link %s: %w", id, err)
	}

	defer lock.Release()

	offset, err := strconv.ParseUint(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		offset = 0
	}

	limit, err := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
	if err != nil || limit > 100 {
		limit = 100
	}

	files, err := s.fs.ListLinkFiles(id)
	if err != nil {
		return fmt.Errorf("failed to get all files for link %s: %w", id, err)
	}

	// TODO: pagination UI elements

	var params template.Params[template.UserViewLinkParams]
	params.Title = "Send me the Data"
	params.Data.Files = view.Files(lock, files, int(offset), int(limit))
	params.Data.Link = view.Link(lock, files)

	lock.Release()

	flashes := flash.GetFlashes(w, r)
	params.SuccessFlash = flashes[flash.SuccessFlash]
	params.ErrorFlash = flashes[flash.ErrorFlash]

	return template.RenderUserViewLink(w, params)
}
