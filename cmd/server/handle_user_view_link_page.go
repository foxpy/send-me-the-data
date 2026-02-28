package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"

	_ "embed"
)

var (
	//go:embed templates/user_view_link.gohtml
	uploadTemplateStr string

	uploadTemplate = template.Must(template.New("").Parse(uploadTemplateStr))
)

func (s *State) handleUserViewLinkPage(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	files, err := s.prepareFilesView(id)
	if errors.Is(err, sql.ErrNoRows) {
		return respond404(w)
	} else if err != nil {
		return err
	}

	successFlash := false
	_, err = r.Cookie("success_flash")
	if err == nil {
		successFlash = true
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "success_flash",
		MaxAge: -1,
	})

	var b bytes.Buffer
	err = uploadTemplate.Execute(&b, map[string]any{
		"Files":        files,
		"SuccessFlash": successFlash,
	})
	if err != nil {
		return fmt.Errorf("failed to render a template: %w", err)
	}

	_, err = io.Copy(w, &b)
	if err != nil {
		return fmt.Errorf("failed to write rendered template: %w", err)
	}

	return nil
}

func (s *State) prepareFilesView(id string) ([]FileView, error) {
	lock, err := s.db.AcquireLinkRLock(id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("failed to acquire read lock on link %s: %w", id, err)
	}

	defer func() {
		_ = lock.Close()
	}()

	files, err := s.getFilesView(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get files view for link %s: %w", id, err)
	}

	return files, nil
}
