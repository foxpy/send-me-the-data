package template

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
	Title        string
	SuccessFlash string
	ErrorFlash   string
	Data         T
}

type UserViewLinkParams struct {
	Link  LinkView
	Files []FileView
}

type AdminViewLinkParams struct {
	Link  LinkView
	Files []FileView
}

type AdminEditLinkParams struct {
	Link LinkView
}

type AdminViewLinksParams struct {
	Links []LinkView
}

type FileView struct {
	Name              string
	UploadedAt        uint64
	Size              string
	AdminDownloadLink string
	UserDownloadLink  string
	DeleteLink        string
}

type LinkView struct {
	Name             string
	CreatedAt        uint64
	TotalFiles       int
	TotalSize        string
	MaxFileSize      string
	MaxFileSizeBytes string
	ViewLink         string
	DeleteLink       string
	EditLink         string
	DownloadZIP      string
	UserDownloadable bool
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

func RenderError(w http.ResponseWriter, code int) error {
	return renderHelper(w, "error.gohtml", fmt.Sprintf("%d %s", code, http.StatusText(code)))
}

func RenderUserViewLink(w http.ResponseWriter, params Params[UserViewLinkParams]) error {
	return renderHelper(w, "user_view_link.gohtml", params)
}

func RenderAdminViewLink(w http.ResponseWriter, params Params[AdminViewLinkParams]) error {
	return renderHelper(w, "admin_view_link.gohtml", params)
}

func RenderAdminEditLink(w http.ResponseWriter, params Params[AdminEditLinkParams]) error {
	return renderHelper(w, "admin_edit_link.gohtml", params)
}

func RenderAdminViewLinks(w http.ResponseWriter, params Params[AdminViewLinksParams]) error {
	return renderHelper(w, "admin_view_links.gohtml", params)
}
