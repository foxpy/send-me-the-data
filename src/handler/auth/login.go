package auth

import (
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

type passwordCheckRequest struct {
	password     string
	passwordHash []byte
	response     chan<- bool
}

func passwordChecker(c <-chan passwordCheckRequest) {
	for pcr := range c {
		err := bcrypt.CompareHashAndPassword(pcr.passwordHash, []byte(pcr.password))
		if err == nil {
			pcr.response <- true
		} else {
			pcr.response <- false
		}
	}
}

func LoginHandler(
	db idb.Database,
	rnd irnd.Random,
	loginURL string,
	successRedirectURL string,
	passwordHashThreads uint,
) http.HandlerFunc {
	if passwordHashThreads == 0 {
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		}
	}

	requestChan := make(chan passwordCheckRequest)
	slog.Info("starting password hash checker", "threads", passwordHashThreads)
	// TODO: graceful shutdown
	for range passwordHashThreads {
		go passwordChecker(requestChan)
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
