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
	"github.com/foxpy/send-me-the-data/src/irnd"
	"golang.org/x/crypto/bcrypt"
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
		http.Redirect(w, r, a.loginURL, http.StatusSeeOther)
		return
	}

	if token.ExpiresAt.Before(time.Now()) {
		err = a.db.DeleteSessionToken(tokenCookie.Value)
		if err != nil {
			slog.Error("failed to delete session token from database", "error", err)
		}
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

type passwordChecker struct{}

type passwordCheckRequest struct {
	password     string
	passwordHash []byte
	response     chan<- bool
}

func (pc passwordChecker) run(c <-chan passwordCheckRequest) {
	for pcr := range c {
		err := bcrypt.CompareHashAndPassword(pcr.passwordHash, []byte(pcr.password))
		if err == nil {
			pcr.response <- true
		} else {
			pcr.response <- false
		}
	}
}

func LoginHandler(db idb.Database, rnd irnd.Random, loginURL, successRedirectURL string) http.HandlerFunc {
	pc := passwordChecker{}
	requestChan := make(chan passwordCheckRequest)
	// TODO: graceful shutdown
	// TODO: configurable parallelism
	for range 2 {
		go pc.run(requestChan)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: check if already authenticated, redirect immediately on success
		// TODO: I will need to implement a log out button first

		username := r.FormValue("user")
		password := r.FormValue("password")
		if len(username) == 0 || len(password) == 0 {
			http.Redirect(w, r, loginURL, http.StatusSeeOther)
			return
		}

		// TODO: maybe remember last time we called this thing and only call it at most every few minutes
		err := db.DeleteOutdatedSessionTokens()
		if err != nil {
			slog.Error("failed to delete outdated session tokens", "error", err)
		}

		passwordHash, err := db.GetAdminPasswordHash(username)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				slog.Error("failed to get password hash from database", "error", err)
			}
			flash.AddFlash(w, flash.ErrorFlash, "Invalid username or password")
			http.Redirect(w, r, loginURL, http.StatusSeeOther)
			return
		}

		responseChan := make(chan bool)
		pcRequest := passwordCheckRequest{
			password:     password,
			passwordHash: passwordHash,
			response:     responseChan,
		}
		requestChan <- pcRequest
		passwordMatches := <-responseChan
		close(responseChan)
		if !passwordMatches {
			flash.AddFlash(w, flash.ErrorFlash, "Invalid username or password")
			http.Redirect(w, r, loginURL, http.StatusSeeOther)
			return
		}

		token := idb.SessionToken{
			Token:     rnd.SessionToken(),
			Username:  username,
			ExpiresAt: time.Now().Add(sessionTokenDuration),
		}
		err = db.CreateSessionToken(token)
		if err != nil {
			slog.Error("failed to write session token to database", "error", err, "username", username)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     sessionTokenCookieName,
			Value:    token.Token,
			Path:     "/",
			Expires:  token.ExpiresAt,
			Secure:   true,                    // prevents the browser from sending this cookie over HTTP
			HttpOnly: true,                    // prevents JS from reading this cookie
			SameSite: http.SameSiteStrictMode, // prevents this cookie from being used by cross-origin requests
		})
		http.Redirect(w, r, successRedirectURL, http.StatusSeeOther)
	}
}
