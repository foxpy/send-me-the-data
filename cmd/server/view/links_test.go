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
				CreatedAt:        mockTime,
				UserDownloadable: false,
				MaxFileSize:      100,
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
					CreatedAt:        mockTimeMilli,
					TotalFiles:       0,
					TotalSize:        "0 bytes",
					MaxFileSize:      "100 bytes",
					MaxFileSizeBytes: "100",
					ViewLink:         "/link/abcd",
					DeleteLink:       "/link/abcd/delete",
					EditLink:         "/link/abcd/edit",
					DownloadZIP:      "/link/abcd/zip",
					UserDownloadable: false,
				},
			},
		},
		{
			desc: "one link with files",
			links: []idb.Link{{
				Name:             "test 1",
				ExternalKey:      "abcd",
				CreatedAt:        mockTime,
				UserDownloadable: true,
				MaxFileSize:      10240,
			}},
			files: []testutil.LinkFiles{
				{
					Name: "abcd",
					Files: []ifs.File{
						{
							Name:    "file 1",
							Size:    100,
							ModTime: mockTime,
						},
						{
							Name:    "file 2",
							Size:    500,
							ModTime: mockTime,
						},
					},
				},
			},
			res: []template.LinkView{
				{
					Name:             "test 1",
					CreatedAt:        mockTimeMilli,
					TotalFiles:       2,
					TotalSize:        "600 bytes",
					MaxFileSize:      "10.00 KiB",
					MaxFileSizeBytes: "10240",
					ViewLink:         "/link/abcd",
					DeleteLink:       "/link/abcd/delete",
					EditLink:         "/link/abcd/edit",
					DownloadZIP:      "/link/abcd/zip",
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
					CreatedAt:        mockTime,
					UserDownloadable: false,
					MaxFileSize:      10240,
				},
				{
					Name:             "test 2",
					ExternalKey:      "bcde",
					CreatedAt:        mockTime,
					UserDownloadable: false,
					MaxFileSize:      100,
				},
			},
			files: []testutil.LinkFiles{
				{
					Name: "abcd",
					Files: []ifs.File{
						{
							Name:    "file 1",
							Size:    100,
							ModTime: mockTime,
						},
						{
							Name:    "file 2",
							Size:    500,
							ModTime: mockTime,
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
					CreatedAt:        mockTimeMilli,
					TotalFiles:       2,
					TotalSize:        "600 bytes",
					MaxFileSize:      "10.00 KiB",
					MaxFileSizeBytes: "10240",
					ViewLink:         "/link/abcd",
					DeleteLink:       "/link/abcd/delete",
					EditLink:         "/link/abcd/edit",
					DownloadZIP:      "/link/abcd/zip",
					UserDownloadable: false,
				},
				{
					Name:             "test 2",
					CreatedAt:        mockTimeMilli,
					TotalFiles:       0,
					TotalSize:        "0 bytes",
					MaxFileSize:      "100 bytes",
					MaxFileSizeBytes: "100",
					ViewLink:         "/link/bcde",
					DeleteLink:       "/link/bcde/delete",
					EditLink:         "/link/bcde/edit",
					DownloadZIP:      "/link/bcde/zip",
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
	for _, tc := range []struct {
		name             string
		linkID           string
		linkName         string
		userDownloadable bool
		maxFileSize      uint64
		files            testutil.LinkFiles
		res              *template.LinkView
	}{
		{
			name:             "no files",
			linkID:           "abcd",
			linkName:         "My Link",
			userDownloadable: false,
			maxFileSize:      4096,
			files: testutil.LinkFiles{
				Name:  "abcd",
				Files: []ifs.File{},
			},
			res: &template.LinkView{
				Name:             "My Link",
				CreatedAt:        mockTimeMilli,
				TotalFiles:       0,
				TotalSize:        "0 bytes",
				MaxFileSize:      "4.00 KiB",
				MaxFileSizeBytes: "4096",
				ViewLink:         "/link/abcd",
				DeleteLink:       "/link/abcd/delete",
				EditLink:         "/link/abcd/edit",
				DownloadZIP:      "/link/abcd/zip",
				UserDownloadable: false,
			},
		},
		{
			name:             "one file",
			linkID:           "abcd",
			linkName:         "My Link",
			userDownloadable: false,
			maxFileSize:      4096,
			files: testutil.LinkFiles{
				Name: "abcd",
				Files: []ifs.File{{
					Name:    "file 1",
					Size:    1024,
					ModTime: mockTime,
				}},
			},
			res: &template.LinkView{
				Name:             "My Link",
				CreatedAt:        mockTimeMilli,
				TotalFiles:       1,
				TotalSize:        "1.00 KiB",
				MaxFileSize:      "4.00 KiB",
				MaxFileSizeBytes: "4096",
				ViewLink:         "/link/abcd",
				DeleteLink:       "/link/abcd/delete",
				EditLink:         "/link/abcd/edit",
				DownloadZIP:      "/link/abcd/zip",
				UserDownloadable: false,
			},
		},
		{
			name:             "many files, user downloadable",
			linkID:           "abcd",
			linkName:         "My Link",
			userDownloadable: true,
			maxFileSize:      4096,
			files: testutil.LinkFiles{
				Name: "abcd",
				Files: []ifs.File{
					{
						Name:    "file 1",
						Size:    1024,
						ModTime: mockTime,
					},
					{
						Name:    "file 2",
						Size:    1024 * 3,
						ModTime: mockTime,
					},
				},
			},
			res: &template.LinkView{
				Name:             "My Link",
				CreatedAt:        mockTimeMilli,
				TotalFiles:       2,
				TotalSize:        "4.00 KiB",
				MaxFileSize:      "4.00 KiB",
				MaxFileSizeBytes: "4096",
				ViewLink:         "/link/abcd",
				DeleteLink:       "/link/abcd/delete",
				EditLink:         "/link/abcd/edit",
				DownloadZIP:      "/link/abcd/zip",
				UserDownloadable: true,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db := mockdb.NewMockDB()
			fs := mockfs.NewMockFS()

			db.SetAcquireLinkRLockResponse(
				tc.linkID,
				tc.linkName,
				mockTime,
				tc.userDownloadable,
				tc.maxFileSize,
			)
			fs.SetListLinkFilesResponse(tc.files.Name, tc.files.Files)
			lock, err := db.AcquireLinkRLock(tc.linkID)
			if err != nil {
				t.Fatal(err)
			}

			linkView, err := Link(lock, fs)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(linkView, tc.res) {
				t.Fatalf("expected %v, got %v", tc.res, linkView)
			}
		})
	}
}
