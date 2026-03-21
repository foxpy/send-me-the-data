package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/foxpy/send-me-the-data/cmd/server/template"
	"github.com/foxpy/send-me-the-data/cmd/server/view"
)

func (s *State) handleUserViewLinkPage(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	title, files, err := s.prepareFilesView(id, false)
	if errors.Is(err, sql.ErrNoRows) {
		return respond404(w)
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

func (s *State) prepareFilesView(id string, forAdmin bool) (string, []template.FileView, error) {
	lock, err := s.db.AcquireLinkRLock(id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil, err
	} else if err != nil {
		return "", nil, fmt.Errorf("failed to acquire read lock on link %s: %w", id, err)
	}

	defer func() {
		_ = lock.Close()
	}()

	renderDownloadLinks := forAdmin || lock.UserDownloadable()
	files, err := view.Files(s.fs, id, renderDownloadLinks)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get files view for link %s: %w", id, err)
	}

	return lock.Name(), files, nil
}
