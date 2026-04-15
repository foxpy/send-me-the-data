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
	m.HandleFunc("GET /{id}", handler.HandleWith500OnError(s.viewLinkPage))
	m.HandleFunc("POST /{id}", handler.HandleWith500OnError(s.upload))
	m.HandleFunc("GET /{id}/{name}", func(w http.ResponseWriter, r *http.Request) {
		if r.PathValue("id") == "static" {
			http.FileServerFS(handler.Static).ServeHTTP(w, r)
		} else {
			handler.HandleWith500OnError(s.downloadFile).ServeHTTP(w, r)
		}
	})
	return handler.WithLogger(m, "user")
}
