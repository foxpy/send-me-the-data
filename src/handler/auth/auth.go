package auth

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/foxpy/send-me-the-data/src/flash"
	"github.com/foxpy/send-me-the-data/src/idb"
)

// TODO: When redirecting to /login, remember where user was going to and save it in a cookie,
//       then after successful authentication, redirect them using this cookie (if it is set).
//       Each successfull authentication should always erase this cookie.

type authenticationMiddleware struct {
	handler  http.Handler
	loginURL string
	db       idb.Database
}

type username struct{}

const sessionTokenCookieName = "session_token"
const sessionTokenDuration = 30 * 24 * time.Hour

func (a *authenticationMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tokenCookie, err := r.Cookie(sessionTokenCookieName)
	if err != nil {
		if err != http.ErrNoCookie {
			slog.Error("failed to obtain session token cookie from user request", "error", err)
		}
		http.Redirect(w, r, a.loginURL, http.StatusSeeOther)
		return
	}

	// TODO: maybe it is a good idea to cache them in memory
	token, err := a.db.GetSessionToken(tokenCookie.Value)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.Error("failed to get session token from database", "error", err)
		}
		// TODO: set cookie: delete session_token
		http.Redirect(w, r, a.loginURL, http.StatusSeeOther)
		return
	}

	if token.ExpiresAt.Before(time.Now()) {
		err = a.db.DeleteSessionToken(tokenCookie.Value)
		if err != nil {
			slog.Error("failed to delete session token from database", "error", err)
		}
		// TODO: set cookie: delete session_token
		flash.AddFlash(w, flash.ErrorFlash, "Session expired")
		http.Redirect(w, r, a.loginURL, http.StatusSeeOther)
		return
	}

	ctx := context.WithValue(r.Context(), username{}, token.Username)
	a.handler.ServeHTTP(w, r.WithContext(ctx))
}

func WithAuthentication(handler http.Handler, loginURL string, db idb.Database) http.Handler {
	return &authenticationMiddleware{handler, loginURL, db}
}

func GetUsername(r *http.Request) string {
	username, _ := r.Context().Value(username{}).(string)
	return username
}
