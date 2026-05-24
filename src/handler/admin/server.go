package admin

import (
	"net/http"

	"github.com/foxpy/send-me-the-data/src/handler"
	"github.com/foxpy/send-me-the-data/src/handler/auth"
	"github.com/foxpy/send-me-the-data/src/idb"
	"github.com/foxpy/send-me-the-data/src/ifs"
	"github.com/foxpy/send-me-the-data/src/irnd"
)

type AdminServer struct {
	db  idb.Database
	fs  ifs.Filesystem
	rnd irnd.Random
}

func NewAdminServer(
	db idb.Database,
	fs ifs.Filesystem,
	rnd irnd.Random,
	passwordHashThreads uint,
) http.Handler {
	s := AdminServer{db, fs, rnd}
	m := http.NewServeMux()

	m.HandleFunc("GET /login", handler.HandleWith500OnError(s.loginPage))
	m.HandleFunc("POST /login", auth.LoginHandler(db, rnd, "/login", "/", passwordHashThreads))

	enableAuthentication := false
	if passwordHashThreads > 0 {
		enableAuthentication = true
	}
	m.Handle("/", s.authenticated(enableAuthentication))

	m.Handle("GET /static/", http.FileServerFS(handler.Static))

	return handler.WithLogger(m, "admin")
}

func (s *AdminServer) authenticated(enableAuthentication bool) http.Handler {
	m := http.NewServeMux()

	m.HandleFunc("GET /{$}", handler.HandleWith500OnError(s.viewLinksPage))

	m.HandleFunc("GET /link/{id}", handler.HandleWith500OnError(s.viewLinkPage))
	// TODO: replace POST with DELETE for delete methods
	m.HandleFunc("POST /link/{id}/delete", handler.HandleWith500OnError(s.deleteLink))
	m.HandleFunc("POST /link/{id}/edit", handler.HandleWith500OnError(s.editLink))
	m.HandleFunc("POST /link", handler.HandleWith500OnError(s.createLink))

	m.HandleFunc("GET /link/{id}/zip", handler.HandleWith500OnError(s.downloadZIP))

	m.HandleFunc("GET /link/{id}/file/{name}", handler.HandleWith500OnError(
		func(w http.ResponseWriter, r *http.Request) error {
			return handler.HandleDownloadFile(w, r, s.db, s.fs, false)
		},
	))
	m.HandleFunc("POST /link/{id}/file/{name}/delete", handler.HandleWith500OnError(s.deleteFile))

	if enableAuthentication {
		m.HandleFunc("POST /logout", auth.LogoutHandler(s.db, "/login"))
		return auth.WithAuthentication(m, "/login", s.db)
	} else {
		return m
	}
}
