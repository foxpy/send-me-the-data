package main

import (
	"fmt"
	"net/http"

	"github.com/foxpy/send-me-the-data/cmd/server/templates"
)

func (s *State) handleAdminViewLinksPage(w http.ResponseWriter, r *http.Request) error {
	links, err := s.getLinksView()
	if err != nil {
		return fmt.Errorf("failed to get links view: %w", err)
	}

	var params templates.Params[templates.AdminViewLinksParams]
	params.Data.Links = links

	_, err = r.Cookie("success_flash")
	if err == nil {
		params.SuccessFlash = "Link created successfully"
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "success_flash",
		MaxAge: -1,
	})

	_, err = r.Cookie("error_flash")
	if err == nil {
		params.ErrorFlash = "Failed to create link"
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "error_flash",
		MaxAge: -1,
	})

	return templates.RenderAdminViewLinks(w, params)
}
