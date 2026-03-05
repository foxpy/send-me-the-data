package main

import (
	"bufio"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/lib/pq"
)

var (
	alphabet     []byte
	alphabetSize big.Int
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

	alphabetSize = *big.NewInt(int64(len(alphabet)))
}

func (s *State) handleAdminCreateLink(w http.ResponseWriter, r *http.Request) error {
	name := r.FormValue("name")
	externalKey, err := generateRandomExternalKey()
	if err != nil {
		return fmt.Errorf("failed to generate a random external key: %w", err)
	}

	userDownloadable := false
	if r.FormValue("user_downloadable") == "on" {
		userDownloadable = true
	}

	err = s.db.CreateLink(name, externalKey, userDownloadable)
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

	http.SetCookie(w, &http.Cookie{
		Name:   "success_flash",
		MaxAge: 60,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func generateRandomExternalKey() (string, error) {
	// Usually we only need one byte per character, but there is a probability
	// we will need more. Assuming an already unlikely case of each call to crypto/rand.Int
	// requiring 2 random bytes instead of 1, buffering 24 bytes in advance guarantees
	// that on Linux we will almost never need more than a single call to getrandom(2).
	r := bufio.NewReaderSize(rand.Reader, 24)

	var result strings.Builder
	for range 12 {
		n, err := rand.Int(r, &alphabetSize)
		if err != nil {
			return "", fmt.Errorf("failed to generate random integer: %w", err)
		}

		result.WriteByte(alphabet[n.Int64()])
	}
	return result.String(), nil
}
