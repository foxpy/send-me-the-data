package admin

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"

	"github.com/lib/pq"
)

var (
	alphabet []byte
)

func init() {
	for i := byte('a'); i <= byte('z'); i++ {
		alphabet = append(alphabet, i)
	}
	for i := byte('A'); i <= byte('Z'); i++ {
		alphabet = append(alphabet, i)
	}
	for i := byte('0'); i <= byte('9'); i++ {
		alphabet = append(alphabet, i)
	}
}

func (s *AdminServer) createLink(w http.ResponseWriter, r *http.Request) error {
	// TODO: check that name is at least not of length 0
	name := r.FormValue("name")
	externalKey := generateRandomExternalKey()

	userDownloadable := false
	if r.FormValue("user_downloadable") == "on" {
		userDownloadable = true
	}

	err := s.db.CreateLink(name, externalKey, userDownloadable)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation" {
		http.SetCookie(w, &http.Cookie{
			Name:   "error_flash",
			MaxAge: 60,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to create link: %w", err)
	}

	// TODO: set flash text in cookie value EVERYWHERE
	http.SetCookie(w, &http.Cookie{
		Name:   "success_flash",
		MaxAge: 60,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func generateRandomExternalKey() string {
	var result [12]byte
	for i := range 12 {
		n := rand.IntN(len(alphabet))
		result[i] = alphabet[n]
	}
	return string(result[:])
}
