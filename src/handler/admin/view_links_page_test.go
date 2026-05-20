package admin

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/foxpy/send-me-the-data/src/flash"
	"github.com/foxpy/send-me-the-data/src/flash/flashtest"
	"github.com/foxpy/send-me-the-data/src/idb"
	"github.com/foxpy/send-me-the-data/src/idb/mockdb"
	"github.com/foxpy/send-me-the-data/src/ifs"
	"github.com/foxpy/send-me-the-data/src/ifs/mockfs"
	"github.com/foxpy/send-me-the-data/src/irnd/mockrnd"
	"github.com/foxpy/send-me-the-data/src/testutil"

	"golang.org/x/net/html"
)

type table struct {
	numrows int
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

func TestViewLinksPage(t *testing.T) {
	for _, tc := range []struct {
		name            string
		links           []idb.Link
		files           []testutil.LinkFiles
		cookies         []*http.Cookie
		expectedCode    int
		expectedTables  []table
		expectedFlashes []flashtest.HTMLFlash
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
			name:  "one link, no flashes",
			links: []idb.Link{mockdb.NewLink("abcdef", "link1", testutil.MockTime, false, false, 0)},
			files: []testutil.LinkFiles{{
				Name: "abcdef",
				Files: []ifs.File{
					{
						Name:    "file 1",
						Size:    1024,
						ModTime: testutil.MockTime,
					},
					{
						Name:    "file 2",
						Size:    10240,
						ModTime: testutil.MockTime,
					},
				},
			}},
			cookies:         nil,
			expectedCode:    http.StatusOK,
			expectedTables:  []table{{1}},
			expectedFlashes: nil,
		},
		{
			name:  "one link, success flash",
			links: []idb.Link{mockdb.NewLink("abcdef", "link1", testutil.MockTime, false, false, 0)},
			files: []testutil.LinkFiles{{
				Name: "abcdef",
				Files: []ifs.File{
					{
						Name:    "file 1",
						Size:    1024,
						ModTime: testutil.MockTime,
					},
					{
						Name:    "file 2",
						Size:    10240,
						ModTime: testutil.MockTime,
					},
				},
			}},
			cookies: []*http.Cookie{{
				Name:  "success_flash",
				Value: "Link created successfully",
			}},
			expectedCode:    http.StatusOK,
			expectedTables:  []table{{1}},
			expectedFlashes: []flashtest.HTMLFlash{{Kind: flash.SuccessFlash, Text: "Link created successfully"}},
		},
		{
			name:  "one link, error flash",
			links: []idb.Link{mockdb.NewLink("abcdef", "link1", testutil.MockTime, false, false, 0)},
			files: []testutil.LinkFiles{{
				Name: "abcdef",
				Files: []ifs.File{
					{
						Name:    "file 1",
						Size:    1024,
						ModTime: testutil.MockTime,
					},
					{
						Name:    "file 2",
						Size:    10240,
						ModTime: testutil.MockTime,
					},
				},
			}},
			cookies: []*http.Cookie{{
				Name:  "error_flash",
				Value: "Failed to create link",
			}},
			expectedCode:    http.StatusOK,
			expectedTables:  []table{{1}},
			expectedFlashes: []flashtest.HTMLFlash{{Kind: flash.ErrorFlash, Text: "Failed to create link"}},
		},
		{
			name: "multile links, no flashes",
			links: []idb.Link{
				mockdb.NewLink("abcdef", "link1", testutil.MockTime, false, false, 0),
				mockdb.NewLink("bcdef", "link2", testutil.MockTime, false, false, 0),
				mockdb.NewLink("cdef", "link3", testutil.MockTime, false, false, 0),
			},
			files: []testutil.LinkFiles{
				{
					Name: "abcdef",
					Files: []ifs.File{
						{
							Name:    "file 1",
							Size:    1024,
							ModTime: testutil.MockTime,
						},
						{
							Name:    "file 2",
							Size:    10240,
							ModTime: testutil.MockTime,
						},
					},
				},
				{
					Name: "bcdef",
					Files: []ifs.File{
						{
							Name:    "file 1",
							Size:    1024,
							ModTime: testutil.MockTime,
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
			fs := mockfs.NewMockFS()
			rnd := mockrnd.NewMockRND()
			h := NewAdminServer(db, fs, rnd)

			defer db.CheckAllExpects()

			// TODO: test different pagination scenarios
			db.MockListLinksResponse(0, 100, tc.links)
			db.MockTotalLinksResponse(uint64(len(tc.links)))
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

			flashes := flashtest.FindAllFlashes(doc)
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
