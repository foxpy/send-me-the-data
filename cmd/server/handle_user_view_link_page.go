package main

import (
	"bytes"
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
	ok, err := s.db.DoesLinkExist(id)
	if err != nil {
		return fmt.Errorf("failed to check if link is published: %w", err)
	}

	if !ok {
		return respond404(w)
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

	files, err := s.getFilesView(id)
	if err != nil {
		return fmt.Errorf("failed to get files view for link %s: %w", id, err)
	}

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
