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

func NewUserServer(db idb.Database, fs ifs.Filesystem) http.Handler {
	s := UserServer{db, fs}
	m := http.NewServeMux()
	// TODO: I want to make user links even shorter, and for that I will have to implement my own mux,
	//       which matches /static/ and /{id} and selects a corresponding handler
	m.HandleFunc("GET /u/{id}", handler.HandleWith500OnError(s.viewLinkPage))
	m.HandleFunc("POST /u/{id}", handler.HandleWith500OnError(s.upload))
	m.HandleFunc("GET /u/{id}/{name}", handler.HandleWith500OnError(s.downloadFile))
	m.Handle("GET /static/", http.FileServerFS(handler.Static))
	return handler.WithLogger(m, "user")
}
