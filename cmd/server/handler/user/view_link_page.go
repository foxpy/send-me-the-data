package user

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/foxpy/send-me-the-data/cmd/server/flash"
	"github.com/foxpy/send-me-the-data/cmd/server/handler"
	"github.com/foxpy/send-me-the-data/cmd/server/idb"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs"
	"github.com/foxpy/send-me-the-data/cmd/server/template"
	"github.com/foxpy/send-me-the-data/cmd/server/view"
)

func (s *UserServer) viewLinkPage(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	title, files, link, err := prepareFilesView(s.db, s.fs, id)
	if errors.Is(err, sql.ErrNoRows) {
		return handler.RespondError(w, http.StatusNotFound)
	} else if err != nil {
		return err
	}

	var params template.Params[template.UserViewLinkParams]
	params.Title = title
	params.Data.Files = files
	params.Data.Link = *link

	flashes := flash.GetFlashes(w, r)
	params.SuccessFlash = flashes[flash.SuccessFlash]
	params.ErrorFlash = flashes[flash.ErrorFlash]

	return template.RenderUserViewLink(w, params)
}

func prepareFilesView(db idb.Database, fs ifs.Filesystem, id string) (string, []template.FileView, *template.LinkView, error) {
	lock, err := db.AcquireLinkRLock(id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil, nil, err
	} else if err != nil {
		return "", nil, nil, fmt.Errorf("failed to acquire read lock on link %s: %w", id, err)
	}

	defer func() {
		_ = lock.Release()
	}()

	files, err := view.Files(fs, lock)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to get files view for link %s: %w", id, err)
	}

	link, err := view.Link(lock, fs)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to get link view for link %s: %w", lock.ExternalKey(), err)
	}

	return fmt.Sprintf("Upload files: %s", lock.Name()), files, link, nil
}
