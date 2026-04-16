package view

import (
	"reflect"
	"testing"

	"github.com/foxpy/send-me-the-data/cmd/server/idb/mockdb"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs/mockfs"
	"github.com/foxpy/send-me-the-data/cmd/server/template"
	"github.com/foxpy/send-me-the-data/cmd/server/testutil"
)

func TestFiles(t *testing.T) {
	for _, tc := range []struct {
		desc             string
		linkID           string
		linkName         string
		userDownloadable bool
		maxFileSize      uint64
		files            []testutil.LinkFiles
		res              []template.FileView
	}{
		{
			desc:             "no files",
			linkID:           "abcd",
			linkName:         "My Link",
			userDownloadable: false,
			maxFileSize:      4096,
			files: []testutil.LinkFiles{{
				Name:  "abcd",
				Files: []ifs.File{},
			}},
			res: []template.FileView{},
		},
		{
			desc:             "one file",
			linkID:           "abcd",
			linkName:         "My Link",
			userDownloadable: false,
			maxFileSize:      4096,
			files: []testutil.LinkFiles{{
				Name: "abcd",
				Files: []ifs.File{{
					Name:    "file 1",
					Size:    1024,
					ModTime: mockTime,
				}},
			}},
			res: []template.FileView{{
				Name:              "file 1",
				UploadedAt:        mockTimeMilli,
				Size:              "1.00 KiB",
				AdminDownloadLink: "/link/abcd/file/file 1",
				UserDownloadLink:  "",
				DeleteLink:        "/link/abcd/file/file 1/delete",
			}},
		},
		{
			desc:             "user downloadable",
			linkID:           "abcd",
			linkName:         "My Link",
			userDownloadable: true,
			maxFileSize:      4096,
			files: []testutil.LinkFiles{{
				Name: "abcd",
				Files: []ifs.File{{
					Name:    "file 1",
					Size:    1024,
					ModTime: mockTime,
				}},
			}},
			res: []template.FileView{{
				Name:              "file 1",
				UploadedAt:        mockTimeMilli,
				Size:              "1.00 KiB",
				AdminDownloadLink: "/link/abcd/file/file 1",
				UserDownloadLink:  "/abcd/file 1",
				DeleteLink:        "/link/abcd/file/file 1/delete",
			}},
		},
		{
			desc:             "many files",
			linkID:           "abcd",
			linkName:         "My Link",
			userDownloadable: false,
			maxFileSize:      4096,
			files: []testutil.LinkFiles{{
				Name: "abcd",
				Files: []ifs.File{
					{
						Name:    "file 1",
						Size:    1024,
						ModTime: mockTime,
					},
					{
						Name:    "file 2",
						Size:    512,
						ModTime: mockTime,
					},
					{
						Name:    "file 3",
						Size:    512,
						ModTime: mockTime,
					},
				},
			}},
			res: []template.FileView{
				{
					Name:              "file 1",
					UploadedAt:        mockTimeMilli,
					Size:              "1.00 KiB",
					AdminDownloadLink: "/link/abcd/file/file 1",
					UserDownloadLink:  "",
					DeleteLink:        "/link/abcd/file/file 1/delete",
				},
				{
					Name:              "file 2",
					UploadedAt:        mockTimeMilli,
					Size:              "512 bytes",
					AdminDownloadLink: "/link/abcd/file/file 2",
					UserDownloadLink:  "",
					DeleteLink:        "/link/abcd/file/file 2/delete",
				},
				{
					Name:              "file 3",
					UploadedAt:        mockTimeMilli,
					Size:              "512 bytes",
					AdminDownloadLink: "/link/abcd/file/file 3",
					UserDownloadLink:  "",
					DeleteLink:        "/link/abcd/file/file 3/delete",
				},
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			db := mockdb.NewMockDB()
			fs := mockfs.NewMockFS()

			db.SetAcquireLinkRLockResponse(tc.linkID, tc.linkName, mockTime, tc.userDownloadable, false, tc.maxFileSize)
			lock, err := db.AcquireLinkRLock(tc.linkID)
			if err != nil {
				t.Fatal(err)
			}

			for _, f := range tc.files {
				fs.SetListLinkFilesResponse(f.Name, f.Files)
			}

			fileViews, err := Files(fs, lock)
			if err != nil {
				t.Fatalf("%s", err)
			}

			if len(fileViews) != len(tc.res) {
				t.Fatalf("expected %d rendered files, got %d", len(tc.res), len(fileViews))
			}

			for i := range fileViews {
				if !reflect.DeepEqual(fileViews[i], tc.res[i]) {
					t.Fatalf("expected %v, got %v", tc.res[i], fileViews[i])
				}
			}
		})
	}
}
