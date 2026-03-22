package admin

import (
	"net/http"

	"github.com/foxpy/send-me-the-data/cmd/server/handler"
	"github.com/foxpy/send-me-the-data/cmd/server/idb"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs"
)

type AdminServer struct {
	db idb.Database
	fs ifs.Filesystem
}

func NewAdminServer(db idb.Database, fs ifs.Filesystem) *http.ServeMux {
	s := AdminServer{db, fs}
	m := http.NewServeMux()

	m.HandleFunc("GET /{$}", handler.HandleWith500OnError(s.viewLinksPage))

	m.HandleFunc("GET /link/{id}", handler.HandleWith500OnError(s.viewLinkPage))
	// TODO: relace POST with DELETE for delete methods
	m.HandleFunc("POST /link/{id}/delete", handler.HandleWith500OnError(s.deleteLink))
	m.HandleFunc("GET /link/{id}/edit", handler.HandleWith500OnError(s.editLinkPage))
	m.HandleFunc("POST /link/{id}/edit", handler.HandleWith500OnError(s.editLink))
	m.HandleFunc("POST /link", handler.HandleWith500OnError(s.createLink))

	m.HandleFunc("GET /link/{id}/file/{name}", handler.HandleWith500OnError(s.downloadFile))
	m.HandleFunc("POST /link/{id}/file/{name}/delete", handler.HandleWith500OnError(s.deleteFile))

	m.Handle("GET /static/", http.FileServerFS(handler.Static))
	return m
}
