package templates

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
)

var (
	//go:embed gohtml/*
	templatesEmbed embed.FS
	templates      = template.Must(template.ParseFS(templatesEmbed, "gohtml/*"))
)

type Params[T any] struct {
	SuccessFlash string
	ErrorFlash   string
	Data         T
}

type UserViewLinkParams struct {
	Files []FileView
}

type AdminViewLinkParams struct {
	Files []FileView
}

type AdminViewLinksParams struct {
	Links []LinkView
}

type FileView struct {
	Name         string
	UploadedAt   string
	Size         string
	DownloadLink string
	DeleteLink   string
}

type LinkView struct {
	Name       string
	CreatedAt  string
	TotalFiles int
	TotalSize  string
	ViewLink   string
	DeleteLink string
}

func renderHelper(w http.ResponseWriter, templateName string, data any) error {
	var b bytes.Buffer
	err := templates.ExecuteTemplate(&b, templateName, data)
	if err != nil {
		return fmt.Errorf("failed to render template %s: %w", templateName, err)
	}

	_, err = io.Copy(w, &b)
	if err != nil {
		slog.Error("failed to write response for rendered template", "template_name", templateName, "error", err)
	}

	return nil
}

func Render404(w http.ResponseWriter) error {
	return renderHelper(w, "404.gohtml", nil)
}

func Render500(w http.ResponseWriter) error {
	return renderHelper(w, "500.gohtml", nil)
}

func RenderUserViewLink(w http.ResponseWriter, params Params[UserViewLinkParams]) error {
	return renderHelper(w, "user_view_link.gohtml", params)
}

func RenderAdminViewLink(w http.ResponseWriter, params Params[AdminViewLinkParams]) error {
	return renderHelper(w, "admin_view_link.gohtml", params)
}

func RenderAdminViewLinks(w http.ResponseWriter, params Params[AdminViewLinksParams]) error {
	return renderHelper(w, "admin_view_links.gohtml", params)
}
