package auth

import (
	"log/slog"
	"net/http"

	"github.com/foxpy/send-me-the-data/src/idb"
)

func LogoutHandler(db idb.Database, loginURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := GetAuth(r)
		if token != nil {
			err := db.DeleteSessionToken(token.Token)
			if err != nil {
				slog.Error("failed to delete session token", "error", err)
			}
		}

		http.SetCookie(w, &http.Cookie{
			Name:   sessionTokenCookieName,
			Path:   "/",
			MaxAge: -1,
		})
		http.Redirect(w, r, loginURL, http.StatusSeeOther)
	}
}
