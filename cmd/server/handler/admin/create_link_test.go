package admin

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/foxpy/send-me-the-data/cmd/server/idb/mockdb"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs/mockfs"
)

func TestCreateLinkEmptyName(t *testing.T) {
	db := mockdb.NewMockDB()
	defer db.CheckAllExpects()
	fs := mockfs.NewMockFS()
	h := NewAdminServer(db, fs)

	db.MockGenerateRandomExternalKeyResponse("abcd")

	// the application doesn't validate link name length, that's the job of the database.
	// in real deployment, postgres would reject such a transaction.
	db.MockExpectedCreateLinkCall("", "abcd", false, 0, func() error {
		return errors.New("mocked postgresql error: link name length constraint validation failure")
	})

	postValues := make(url.Values)
	postValues.Add("max_file_size", "0")
	req := httptest.NewRequest("POST", "/link", bytes.NewReader([]byte(postValues.Encode())))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf(
			"expected status code %s, got %s",
			http.StatusText(http.StatusInternalServerError),
			http.StatusText(resp.StatusCode),
		)
	}
}

func TestCreateLink(t *testing.T) {
	db := mockdb.NewMockDB()
	defer db.CheckAllExpects()
	fs := mockfs.NewMockFS()
	h := NewAdminServer(db, fs)

	db.MockGenerateRandomExternalKeyResponse("abcd")
	db.MockExpectedCreateLinkCall("My Link", "abcd", false, 9000, nil)

	postValues := make(url.Values)
	postValues.Add("name", "My Link")
	postValues.Add("max_file_size", "9000")
	req := httptest.NewRequest("POST", "/link", bytes.NewReader([]byte(postValues.Encode())))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther {
		t.Fatalf(
			"expected status code %s, got %s",
			http.StatusText(http.StatusSeeOther),
			http.StatusText(resp.StatusCode),
		)
	}

	expectedLocation := "/"
	location := resp.Header.Get("Location")
	if location != expectedLocation {
		t.Fatalf("expected redirect to %s, got %s", expectedLocation, location)
	}
}
