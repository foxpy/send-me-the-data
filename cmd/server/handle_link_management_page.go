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
	//go:embed templates/link_management.gohtml
	linkManagementStr string

	linkManagementTemplate = template.Must(template.New("").Parse(linkManagementStr))
)

func (s *State) handleLinkManagementPage(w http.ResponseWriter, r *http.Request) error {
	links, err := listLinks(s.prefix, s.db)
	if err != nil {
		return fmt.Errorf("failed to list links: %w", err)
	}

	var b bytes.Buffer
	err = linkManagementTemplate.Execute(&b, links)
	if err != nil {
		return fmt.Errorf("failed to render a template: %w", err)
	}

	_, err = io.Copy(w, &b)
	if err != nil {
		return fmt.Errorf("failed to write rendered template: %w", err)
	}

	return nil

	// TODO: Create link form
}
