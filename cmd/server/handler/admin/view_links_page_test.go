package admin

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/foxpy/send-me-the-data/cmd/server/idb"
	"github.com/foxpy/send-me-the-data/cmd/server/idb/mockdb"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs/mockfs"
	"github.com/foxpy/send-me-the-data/cmd/server/testutil"

	"golang.org/x/net/html"
)

type table struct {
	numrows int
}

type flash struct {
	kind string
	text string
}

func tableRows(table *html.Node) int {
	var tbody *html.Node
	for n := range table.ChildNodes() {
		if n.Type == html.ElementNode && n.Data == "tbody" {
			tbody = n
			break
		}
	}
	if tbody == nil {
		panic("table without tbody")
	}

	rows := 0
	for n := range tbody.ChildNodes() {
		if n.Type == html.ElementNode && n.Data == "tr" {
			rows++
		}
	}
	return rows
}

func findAllTables(doc *html.Node) (tables []*html.Node) {
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.Data == "table" {
			tables = append(tables, n)
		}
	}
	return
}

func findAllFlashes(doc *html.Node) (flashes []flash) {
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode {
			for _, a := range n.Attr {
				if a.Key == "class" {
					kind := a.Val
					if kind != "success_flash" && kind != "error_flash" {
						continue
					}
					flashes = append(flashes, flash{
						kind: kind,
						text: strings.TrimSpace(n.FirstChild.Data),
					})
				}
			}
		}
	}
	return
}

func TestViewLinksPage(t *testing.T) {
	for _, tc := range []struct {
		name            string
		links           []idb.Link
		files           []testutil.LinkFiles
		cookies         []*http.Cookie
		expectedCode    int
		expectedTables  []table
		expectedFlashes []flash
	}{
		{
			name:            "no links, no flashes",
			links:           []idb.Link{},
			files:           nil,
			cookies:         nil,
			expectedCode:    http.StatusOK,
			expectedTables:  []table{{0}},
			expectedFlashes: nil,
		},
		{
			name: "one link, no flashes",
			links: []idb.Link{{
				Name:        "link1",
				ExternalKey: "abcdef",
				CreatedAt:   time.UnixMicro(0).UTC(),
			}},
			files: []testutil.LinkFiles{{
				Name: "abcdef",
				Files: []ifs.File{
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
				},
			}},
			cookies:         nil,
			expectedCode:    http.StatusOK,
			expectedTables:  []table{{1}},
			expectedFlashes: nil,
		},
		{
			name: "one link, success flash",
			links: []idb.Link{{
				Name:        "link1",
				ExternalKey: "abcdef",
				CreatedAt:   time.UnixMicro(0).UTC(),
			}},
			files: []testutil.LinkFiles{{
				Name: "abcdef",
				Files: []ifs.File{
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
				},
			}},
			cookies: []*http.Cookie{{
				Name: "success_flash",
			}},
			expectedCode:    http.StatusOK,
			expectedTables:  []table{{1}},
			expectedFlashes: []flash{{"success_flash", "Link created successfully"}},
		},
		{
			name: "one link, error flash",
			links: []idb.Link{{
				Name:        "link1",
				ExternalKey: "abcdef",
				CreatedAt:   time.UnixMicro(0).UTC(),
			}},
			files: []testutil.LinkFiles{{
				Name: "abcdef",
				Files: []ifs.File{
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
				},
			}},
			cookies: []*http.Cookie{{
				Name: "error_flash",
			}},
			expectedCode:    http.StatusOK,
			expectedTables:  []table{{1}},
			expectedFlashes: []flash{{"error_flash", "Failed to create link"}},
		},
		{
			name: "multile links, no flashes",
			links: []idb.Link{
				{
					Name:        "link1",
					ExternalKey: "abcdef",
					CreatedAt:   time.UnixMicro(0).UTC(),
				},
				{
					Name:        "link2",
					ExternalKey: "bcdef",
					CreatedAt:   time.UnixMicro(0).UTC(),
				},
				{
					Name:        "link3",
					ExternalKey: "cdef",
					CreatedAt:   time.UnixMicro(0).UTC(),
				},
			},
			files: []testutil.LinkFiles{
				{
					Name: "abcdef",
					Files: []ifs.File{
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
					},
				},
				{
					Name: "bcdef",
					Files: []ifs.File{
						{
							Name:    "file 1",
							Size:    1024,
							ModTime: time.UnixMicro(0).UTC(),
						},
					},
				},
				{
					Name:  "cdef",
					Files: []ifs.File{},
				},
			},
			cookies:         nil,
			expectedCode:    http.StatusOK,
			expectedTables:  []table{{3}},
			expectedFlashes: nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db := mockdb.NewMockDB()
			defer db.CheckAllExpects()
			fs := mockfs.NewMockFS()
			h := NewAdminServer(db, fs)

			db.MockAllLinksResponse(tc.links)
			for _, f := range tc.files {
				fs.SetListLinkFilesResponse(f.Name, f.Files)
			}

			req := httptest.NewRequest("GET", "/", nil)
			for _, cookie := range tc.cookies {
				req.AddCookie(cookie)
			}
			w := httptest.NewRecorder()

			h.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedCode {
				t.Fatalf(
					"Expected status code %s, got %s",
					http.StatusText(tc.expectedCode),
					http.StatusText(resp.StatusCode),
				)
			}

			doc, err := html.Parse(resp.Body)
			if err != nil {
				t.Error(err)
			}

			tables := findAllTables(doc)
			if len(tables) != len(tc.expectedTables) {
				t.Fatalf("expected %d tables, got %d", len(tc.expectedTables), len(tables))
			}

			for i := range tables {
				numrows := tableRows(tables[i])
				if numrows != tc.expectedTables[i].numrows {
					t.Fatalf("expected %d rows, got %d", tc.expectedTables[i].numrows, numrows)
				}
			}

			flashes := findAllFlashes(doc)
			if len(flashes) != len(tc.expectedFlashes) {
				t.Fatalf("expected %d flashes, got %d", len(tc.expectedFlashes), len(flashes))
			}

			for i := range flashes {
				if !reflect.DeepEqual(flashes[i], tc.expectedFlashes[i]) {
					t.Fatalf("expected %v, got %v", tc.expectedFlashes[i], flashes[i])
				}
			}
		})
	}
}
