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
	//go:embed templates/admin_view_links.gohtml
	viewLinksStr string

	viewLinksTemplate = template.Must(template.New("").Parse(viewLinksStr))
)

func (s *State) handleAdminViewLinksPage(w http.ResponseWriter, r *http.Request) error {
	links, err := s.getLinksView()
	if err != nil {
		return fmt.Errorf("failed to get links view: %w", err)
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

	errorFlash := false
	_, err = r.Cookie("error_flash")
	if err == nil {
		errorFlash = true
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "error_flash",
		MaxAge: -1,
	})

	var b bytes.Buffer
	err = viewLinksTemplate.Execute(&b, map[string]any{
		"Links":        links,
		"SuccessFlash": successFlash,
		"ErrorFlash":   errorFlash,
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
