package view

import (
	"reflect"
	"testing"
	"time"

	"github.com/foxpy/send-me-the-data/cmd/server/idb"
	"github.com/foxpy/send-me-the-data/cmd/server/idb/mockdb"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs/mockfs"
	"github.com/foxpy/send-me-the-data/cmd/server/templates"
)

func TestLinks(t *testing.T) {
	for _, tc := range []struct {
		desc  string
		links []idb.Link
		files []linkFiles
		res   []templates.LinkView
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
				Name:        "test 1",
				ExternalKey: "abcd",
				CreatedAt:   time.UnixMicro(0).UTC(),
			}},
			files: []linkFiles{
				{
					name:  "abcd",
					files: []ifs.File{},
				},
			},
			res: []templates.LinkView{
				{
					Name:       "test 1",
					CreatedAt:  "Jan  1 00:00:00",
					TotalFiles: 0,
					TotalSize:  "0.00 bytes",
					ViewLink:   "/link/abcd",
					DeleteLink: "/link/abcd/delete",
				},
			},
		},
		{
			desc: "one link with files",
			links: []idb.Link{{
				Name:        "test 1",
				ExternalKey: "abcd",
				CreatedAt:   time.UnixMicro(0).UTC(),
			}},
			files: []linkFiles{
				{
					name: "abcd",
					files: []ifs.File{
						{
							Name:    "file 1",
							Size:    100,
							ModTime: time.UnixMicro(0).UTC(),
						},
						{
							Name:    "file 2",
							Size:    500,
							ModTime: time.UnixMicro(0).UTC(),
						},
					},
				},
			},
			res: []templates.LinkView{
				{
					Name:       "test 1",
					CreatedAt:  "Jan  1 00:00:00",
					TotalFiles: 2,
					TotalSize:  "600.00 bytes",
					ViewLink:   "/link/abcd",
					DeleteLink: "/link/abcd/delete",
				},
			},
		},
		{
			desc: "one link with files, one link without",
			links: []idb.Link{
				{
					Name:        "test 1",
					ExternalKey: "abcd",
					CreatedAt:   time.UnixMicro(0).UTC(),
				},
				{
					Name:        "test 2",
					ExternalKey: "bcde",
					CreatedAt:   time.UnixMicro(0).UTC(),
				},
			},
			files: []linkFiles{
				{
					name: "abcd",
					files: []ifs.File{
						{
							Name:    "file 1",
							Size:    100,
							ModTime: time.UnixMicro(0).UTC(),
						},
						{
							Name:    "file 2",
							Size:    500,
							ModTime: time.UnixMicro(0).UTC(),
						},
					},
				},
				{
					name:  "bcde",
					files: []ifs.File{},
				},
			},
			res: []templates.LinkView{
				{
					Name:       "test 1",
					CreatedAt:  "Jan  1 00:00:00",
					TotalFiles: 2,
					TotalSize:  "600.00 bytes",
					ViewLink:   "/link/abcd",
					DeleteLink: "/link/abcd/delete",
				},
				{
					Name:       "test 2",
					CreatedAt:  "Jan  1 00:00:00",
					TotalFiles: 0,
					TotalSize:  "0.00 bytes",
					ViewLink:   "/link/bcde",
					DeleteLink: "/link/bcde/delete",
				},
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			db := mockdb.NewMockDB()
			fs := mockfs.NewMockFS()

			db.SetAllLinksResponse(tc.links)
			for _, f := range tc.files {
				fs.SetListLinkFilesResponse(f.name, f.files)
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
