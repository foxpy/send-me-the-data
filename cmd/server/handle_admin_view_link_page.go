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
	//go:embed templates/admin_view_link.gohtml
	viewTemplateStr string

	viewTemplate = template.Must(template.New("").Parse(viewTemplateStr))
)

func (s *State) handleAdminViewLinkPage(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	files, err := s.prepareFilesView(id, true)
	if errors.Is(err, sql.ErrNoRows) {
		return respond404(w)
	} else if err != nil {
		return err
	}

	var b bytes.Buffer
	err = viewTemplate.Execute(&b, files)
	if err != nil {
		return fmt.Errorf("failed to render a template: %w", err)
	}

	_, err = io.Copy(w, &b)
	if err != nil {
		return fmt.Errorf("failed to write rendered template: %w", err)
	}

	return nil
}
