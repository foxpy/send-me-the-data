package user

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

func (s *UserServer) viewLinkPage(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	title, files, err := prepareFilesView(s.db, s.fs, id)
	if errors.Is(err, sql.ErrNoRows) {
		return handler.Respond404(w)
	} else if err != nil {
		return err
	}

	var params template.Params[template.UserViewLinkParams]
	params.Title = title
	params.Data.Files = files

	_, err = r.Cookie("success_flash")
	if err == nil {
		params.SuccessFlash = "File uploaded successfully"
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "success_flash",
		MaxAge: -1,
	})

	return template.RenderUserViewLink(w, params)
}

func prepareFilesView(db idb.Database, fs ifs.Filesystem, id string) (string, []template.FileView, error) {
	lock, err := db.AcquireLinkRLock(id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil, err
	} else if err != nil {
		return "", nil, fmt.Errorf("failed to acquire read lock on link %s: %w", id, err)
	}

	defer func() {
		_ = lock.Release()
	}()

	files, err := view.Files(fs, id, lock.UserDownloadable())
	if err != nil {
		return "", nil, fmt.Errorf("failed to get files view for link %s: %w", id, err)
	}

	return fmt.Sprintf("Upload files: %s", lock.Name()), files, nil
}
