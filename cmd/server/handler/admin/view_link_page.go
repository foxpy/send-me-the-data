package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/foxpy/send-me-the-data/cmd/server/handler"
	"github.com/foxpy/send-me-the-data/cmd/server/idb"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs"
	"github.com/foxpy/send-me-the-data/cmd/server/template"
	"github.com/foxpy/send-me-the-data/cmd/server/view"
)

func (s *AdminServer) viewLinkPage(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	title, files, link, err := prepareFilesView(s.db, s.fs, id)
	if errors.Is(err, sql.ErrNoRows) {
		return handler.Respond404(w)
	} else if err != nil {
		return err
	}

	var params template.Params[template.AdminViewLinkParams]
	params.Title = title
	params.Data.Files = files
	params.Data.Link = *link

	_, err = r.Cookie("success_flash")
	if err == nil {
		params.SuccessFlash = "Link updated successfully"
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "success_flash",
		Path:   "/",
		MaxAge: -1,
	})

	_, err = r.Cookie("error_flash")
	if err == nil {
		params.ErrorFlash = "Failed to update link"
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "error_flash",
		Path:   "/",
		MaxAge: -1,
	})

	return template.RenderAdminViewLink(w, params)
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

	files, err := view.Files(fs, id, true)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to get files view for link %s: %w", id, err)
	}

	link, err := view.Link(lock, fs)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to get %s link view: %w", id, err)
	}

	return fmt.Sprintf("Download files: %s", lock.Name()), files, link, nil
}
