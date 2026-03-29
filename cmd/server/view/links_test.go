package view

import (
	"reflect"
	"testing"

	"github.com/foxpy/send-me-the-data/cmd/server/idb"
	"github.com/foxpy/send-me-the-data/cmd/server/idb/mockdb"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs/mockfs"
	"github.com/foxpy/send-me-the-data/cmd/server/template"
	"github.com/foxpy/send-me-the-data/cmd/server/testutil"
)

func TestLinks(t *testing.T) {
	for _, tc := range []struct {
		desc  string
		links []idb.Link
		files []testutil.LinkFiles
		res   []template.LinkView
	}{
		{
			desc:  "zero links",
			links: []idb.Link{},
			files: nil,
			res:   nil,
		},
		{
			desc: "one link without files",
			links: []idb.Link{{
				Name:             "test 1",
				ExternalKey:      "abcd",
				CreatedAt:        zeroTime,
				UserDownloadable: false,
			}},
			files: []testutil.LinkFiles{
				{
					Name:  "abcd",
					Files: []ifs.File{},
				},
			},
			res: []template.LinkView{
				{
					Name:             "test 1",
					CreatedAt:        "Jan 1 00:00:00 UTC 1970",
					TotalFiles:       0,
					TotalSize:        "0 bytes",
					ViewLink:         "/link/abcd",
					DeleteLink:       "/link/abcd/delete",
					EditLink:         "/link/abcd/edit",
					UserDownloadable: false,
				},
			},
		},
		{
			desc: "one link with files",
			links: []idb.Link{{
				Name:             "test 1",
				ExternalKey:      "abcd",
				CreatedAt:        zeroTime,
				UserDownloadable: true,
			}},
			files: []testutil.LinkFiles{
				{
					Name: "abcd",
					Files: []ifs.File{
						{
							Name:    "file 1",
							Size:    100,
							ModTime: zeroTime,
						},
						{
							Name:    "file 2",
							Size:    500,
							ModTime: zeroTime,
						},
					},
				},
			},
			res: []template.LinkView{
				{
					Name:             "test 1",
					CreatedAt:        "Jan 1 00:00:00 UTC 1970",
					TotalFiles:       2,
					TotalSize:        "600 bytes",
					ViewLink:         "/link/abcd",
					DeleteLink:       "/link/abcd/delete",
					EditLink:         "/link/abcd/edit",
					UserDownloadable: true,
				},
			},
		},
		{
			desc: "one link with files, one link without",
			links: []idb.Link{
				{
					Name:             "test 1",
					ExternalKey:      "abcd",
					CreatedAt:        zeroTime,
					UserDownloadable: false,
				},
				{
					Name:             "test 2",
					ExternalKey:      "bcde",
					CreatedAt:        zeroTime,
					UserDownloadable: false,
				},
			},
			files: []testutil.LinkFiles{
				{
					Name: "abcd",
					Files: []ifs.File{
						{
							Name:    "file 1",
							Size:    100,
							ModTime: zeroTime,
						},
						{
							Name:    "file 2",
							Size:    500,
							ModTime: zeroTime,
						},
					},
				},
				{
					Name:  "bcde",
					Files: []ifs.File{},
				},
			},
			res: []template.LinkView{
				{
					Name:             "test 1",
					CreatedAt:        "Jan 1 00:00:00 UTC 1970",
					TotalFiles:       2,
					TotalSize:        "600 bytes",
					ViewLink:         "/link/abcd",
					DeleteLink:       "/link/abcd/delete",
					EditLink:         "/link/abcd/edit",
					UserDownloadable: false,
				},
				{
					Name:             "test 2",
					CreatedAt:        "Jan 1 00:00:00 UTC 1970",
					TotalFiles:       0,
					TotalSize:        "0 bytes",
					ViewLink:         "/link/bcde",
					DeleteLink:       "/link/bcde/delete",
					EditLink:         "/link/bcde/edit",
					UserDownloadable: false,
				},
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			db := mockdb.NewMockDB()
			defer db.CheckAllExpects()
			fs := mockfs.NewMockFS()

			db.MockAllLinksResponse(tc.links)
			for _, f := range tc.files {
				fs.SetListLinkFilesResponse(f.Name, f.Files)
			}

			linkViews, err := Links(db, fs)
			if err != nil {
				t.Fatalf("%s", err)
			}

			if len(linkViews) != len(tc.res) {
				t.Fatalf("expected %d rendered links, got %d", len(tc.res), len(linkViews))
			}

			for i := range linkViews {
				if !reflect.DeepEqual(linkViews[i], tc.res[i]) {
					t.Fatalf("expected %v, got %v", tc.res[i], linkViews[i])
				}
			}
		})
	}
}

func TestLink(t *testing.T) {
	// TODO
}
