package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	main "github.com/foxpy/send-me-the-data/cmd/server"
	"github.com/foxpy/send-me-the-data/cmd/server/idb"
	"github.com/foxpy/send-me-the-data/cmd/server/idb/mockdb"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs/mockfs"

	"golang.org/x/net/html"
)

func findAllTables(doc *html.Node) (tables []*html.Node) {
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.Data == "table" {
			tables = append(tables, n)
		}
	}
	return
}

func findAllFlashes(doc *html.Node, flashClass string) (flashes []*html.Node) {
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode {
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == flashClass {
					flashes = append(flashes, n)
				}
			}
		}
	}
	return
}

func TestHandleAdminViewLinksPage(t *testing.T) {
	db := mockdb.NewMockDB()
	fs := mockfs.NewMockFS()
	s := main.NewStateFromParts(db, fs)
	h := main.AdminServer(s)

	db.SetAllLinksResponse([]idb.Link{{
		Name:        "link1",
		ExternalKey: "abcdef",
		CreatedAt:   time.UnixMicro(0).UTC(),
	}})
	fs.SetListLinkFilesResponse("abcdef", []ifs.File{
		{
			Name:    "file 1",
			Size:    1024,
			ModTime: time.UnixMicro(0).UTC(),
		},
		{
			Name:    "file 2",
			Size:    10240,
			ModTime: time.UnixMicro(0).UTC(),
		},
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatal("Expected status code OK")
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		t.Fatalf("%s", err)
	}

	tables := findAllTables(doc)
	if len(tables) != 1 {
		t.Fatal("expected strictly 1 table")
	}

	errorFlashes := findAllFlashes(doc, "error_flash")
	if len(errorFlashes) != 0 {
		t.Fatal("no error flashes expected")
	}

	successFlashes := findAllFlashes(doc, "success_flash")
	if len(successFlashes) != 0 {
		t.Fatal("no success flashes expected")
	}
}
