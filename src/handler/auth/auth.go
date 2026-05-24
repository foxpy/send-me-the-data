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

const sessionTokenCookieName = "session_token"
const sessionTokenDuration = 30 * 24 * time.Hour

type authenticationMiddleware struct {
	handler  http.Handler
	loginURL string
	db       idb.Database
}

type authKey struct{}

type authenticationResult struct {
	token     *idb.SessionToken
	userError string
	err       error
}

func authenticationData(db idb.Database, r *http.Request) authenticationResult {
	tokenCookie, err := r.Cookie(sessionTokenCookieName)
	if err != nil {
		if err != http.ErrNoCookie {
			slog.Error("failed to obtain session token cookie from user request", "error", err)
		}
		return authenticationResult{nil, "", err}
	}

	// TODO: maybe it is a good idea to cache them in memory
	token, err := db.GetSessionToken(tokenCookie.Value)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			slog.Error("failed to get session token from database", "error", err)
		}
		return authenticationResult{nil, "Session expired", err}
	}

	if token.ExpiresAt.Before(time.Now()) {
		err = db.DeleteSessionToken(token.Token)
		if err != nil {
			slog.Error("failed to delete session token from database", "error", err)
		}
		return authenticationResult{nil, "Session expired", err}
	}

	return authenticationResult{token, "", nil}
}

func (a *authenticationMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := authenticationData(a.db, r)
	if res.err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:   sessionTokenCookieName,
			Path:   "/",
			MaxAge: -1,
		})
		if len(res.userError) > 0 {
			flash.AddFlash(w, flash.ErrorFlash, res.userError)
		}
		http.Redirect(w, r, a.loginURL, http.StatusSeeOther)
		return
	}

	ctx := context.WithValue(r.Context(), authKey{}, res.token)
	a.handler.ServeHTTP(w, r.WithContext(ctx))
}

func WithAuthentication(handler http.Handler, loginURL string, db idb.Database) http.Handler {
	return &authenticationMiddleware{handler, loginURL, db}
}

func GetAuth(r *http.Request) *idb.SessionToken {
	token, _ := r.Context().Value(authKey{}).(*idb.SessionToken)
	if token == nil {
		return &idb.SessionToken{}
	}

	return token
}
