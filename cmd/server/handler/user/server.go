package user

import (
	"net/http"

	"github.com/foxpy/send-me-the-data/cmd/server/handler"
	"github.com/foxpy/send-me-the-data/cmd/server/idb"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs"
)

type UserServer struct {
	db idb.Database
	fs ifs.Filesystem
}

func NewUserServer(db idb.Database, fs ifs.Filesystem) *http.ServeMux {
	s := UserServer{db, fs}
	m := http.NewServeMux()
	// TODO: I want user links to be as short as possible. Ideally, each link should look like this:
	//       Link: /{id}
	//       File: /{id}/{name}
	m.HandleFunc("GET /u/{id}", handler.HandleWith500OnError(s.viewLinkPage))
	m.HandleFunc("POST /u/{id}", handler.HandleWith500OnError(s.upload))
	m.HandleFunc("GET /link/{id}/file/{name}", handler.HandleWith500OnError(s.downloadFile))
	m.Handle("GET /static/", http.FileServerFS(handler.Static))
	return m
}
