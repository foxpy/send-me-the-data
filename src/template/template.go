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
	Pages []PaginationView
}

type AdminViewLinkParams struct {
	Link  LinkView
	Files []FileView
	Pages []PaginationView
}

type AdminViewLinksParams struct {
	Links []LinkView
	Pages []PaginationView
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
	MaxFileSizeBytes uint64
	ViewLink         string
	DeleteLink       string
	EditLink         string
	DownloadZIP      string
	UserDownloadable bool
	UploadEnabled    bool
}

type PaginationView struct {
	Number        uint
	Link          string
	IsPlaceholder bool
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

func RenderAdminViewLinks(w http.ResponseWriter, params Params[AdminViewLinksParams]) error {
	return renderHelper(w, "admin_view_links.gohtml", params)
}

func RenderAdminLogin(w http.ResponseWriter, params Params[struct{}]) error {
	return renderHelper(w, "admin_login.gohtml", params)
}
