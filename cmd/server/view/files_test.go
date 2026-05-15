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

// TODO: more tests with more links to check that all is sorted correctly
// TODO: replace mockTimeMilli with mockTimeNano

func TestFiles(t *testing.T) {
	for _, tc := range []struct {
		desc             string
		linkID           string
		linkName         string
		userDownloadable bool
		maxFileSize      uint64
		files            []testutil.LinkFiles
		offset, limit    uint
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
			offset: 0,
			limit:  100,
			res:    []template.FileView{},
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
					ModTime: testutil.MockTime,
				}},
			}},
			offset: 0,
			limit:  100,
			res: []template.FileView{{
				Name:              "file 1",
				UploadedAt:        testutil.MockTimeMilli,
				Size:              "1 KiB",
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
					ModTime: testutil.MockTime,
				}},
			}},
			offset: 0,
			limit:  100,
			res: []template.FileView{{
				Name:              "file 1",
				UploadedAt:        testutil.MockTimeMilli,
				Size:              "1 KiB",
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
						ModTime: testutil.MockTime.Add(1),
					},
					{
						Name:    "file 3",
						Size:    512,
						ModTime: testutil.MockTime.Add(3),
					},
					{
						Name:    "file 2",
						Size:    512,
						ModTime: testutil.MockTime.Add(2),
					},
				},
			}},
			offset: 0,
			limit:  100,
			res: []template.FileView{
				{
					Name:              "file 3",
					UploadedAt:        testutil.MockTimeMilli,
					Size:              "512 bytes",
					AdminDownloadLink: "/link/abcd/file/file 3",
					UserDownloadLink:  "",
					DeleteLink:        "/link/abcd/file/file 3/delete",
				},
				{
					Name:              "file 2",
					UploadedAt:        testutil.MockTimeMilli,
					Size:              "512 bytes",
					AdminDownloadLink: "/link/abcd/file/file 2",
					UserDownloadLink:  "",
					DeleteLink:        "/link/abcd/file/file 2/delete",
				},
				{
					Name:              "file 1",
					UploadedAt:        testutil.MockTimeMilli,
					Size:              "1 KiB",
					AdminDownloadLink: "/link/abcd/file/file 1",
					UserDownloadLink:  "",
					DeleteLink:        "/link/abcd/file/file 1/delete",
				},
			},
		},
		{
			desc:             "3 files, limit 2",
			linkID:           "abcd",
			linkName:         "My Link",
			userDownloadable: false,
			maxFileSize:      4096,
			files: []testutil.LinkFiles{{
				Name: "abcd",
				Files: []ifs.File{
					{
						Name:    "file 3",
						Size:    512,
						ModTime: testutil.MockTime.Add(3),
					},
					{
						Name:    "file 1",
						Size:    1024,
						ModTime: testutil.MockTime.Add(1),
					},
					{
						Name:    "file 2",
						Size:    512,
						ModTime: testutil.MockTime.Add(2),
					},
				},
			}},
			offset: 0,
			limit:  2,
			res: []template.FileView{
				{
					Name:              "file 3",
					UploadedAt:        testutil.MockTimeMilli,
					Size:              "512 bytes",
					AdminDownloadLink: "/link/abcd/file/file 3",
					UserDownloadLink:  "",
					DeleteLink:        "/link/abcd/file/file 3/delete",
				},
				{
					Name:              "file 2",
					UploadedAt:        testutil.MockTimeMilli,
					Size:              "512 bytes",
					AdminDownloadLink: "/link/abcd/file/file 2",
					UserDownloadLink:  "",
					DeleteLink:        "/link/abcd/file/file 2/delete",
				},
			},
		},
		{
			desc:             "3 files, offset 1",
			linkID:           "abcd",
			linkName:         "My Link",
			userDownloadable: false,
			maxFileSize:      4096,
			files: []testutil.LinkFiles{{
				Name: "abcd",
				Files: []ifs.File{
					{
						Name:    "file 2",
						Size:    512,
						ModTime: testutil.MockTime.Add(2),
					},
					{
						Name:    "file 1",
						Size:    1024,
						ModTime: testutil.MockTime.Add(1),
					},
					{
						Name:    "file 3",
						Size:    512,
						ModTime: testutil.MockTime.Add(3),
					},
				},
			}},
			offset: 1,
			limit:  100,
			res: []template.FileView{
				{
					Name:              "file 2",
					UploadedAt:        testutil.MockTimeMilli,
					Size:              "512 bytes",
					AdminDownloadLink: "/link/abcd/file/file 2",
					UserDownloadLink:  "",
					DeleteLink:        "/link/abcd/file/file 2/delete",
				},
				{
					Name:              "file 1",
					UploadedAt:        testutil.MockTimeMilli,
					Size:              "1 KiB",
					AdminDownloadLink: "/link/abcd/file/file 1",
					UserDownloadLink:  "",
					DeleteLink:        "/link/abcd/file/file 1/delete",
				},
			},
		},
		{
			desc:             "3 files, offset 1, limit 1",
			linkID:           "abcd",
			linkName:         "My Link",
			userDownloadable: false,
			maxFileSize:      4096,
			files: []testutil.LinkFiles{{
				Name: "abcd",
				Files: []ifs.File{
					{
						Name:    "file 2",
						Size:    512,
						ModTime: testutil.MockTime.Add(2),
					},
					{
						Name:    "file 3",
						Size:    512,
						ModTime: testutil.MockTime.Add(3),
					},
					{
						Name:    "file 1",
						Size:    1024,
						ModTime: testutil.MockTime.Add(1),
					},
				},
			}},
			offset: 1,
			limit:  1,
			res: []template.FileView{
				{
					Name:              "file 2",
					UploadedAt:        testutil.MockTimeMilli,
					Size:              "512 bytes",
					AdminDownloadLink: "/link/abcd/file/file 2",
					UserDownloadLink:  "",
					DeleteLink:        "/link/abcd/file/file 2/delete",
				},
			},
		},
		{
			desc:             "3 files, offset 5",
			linkID:           "abcd",
			linkName:         "My Link",
			userDownloadable: false,
			maxFileSize:      4096,
			files: []testutil.LinkFiles{{
				Name: "abcd",
				Files: []ifs.File{
					{
						Name:    "file 3",
						Size:    512,
						ModTime: testutil.MockTime.Add(3),
					},
					{
						Name:    "file 2",
						Size:    512,
						ModTime: testutil.MockTime.Add(2),
					},
					{
						Name:    "file 1",
						Size:    1024,
						ModTime: testutil.MockTime.Add(1),
					},
				},
			}},
			offset: 5,
			limit:  100,
			res:    []template.FileView{},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			db := mockdb.NewMockDB()
			fs := mockfs.NewMockFS()

			db.SetAcquireLinkRLockResponse(tc.linkID, tc.linkName, testutil.MockTime, tc.userDownloadable, false, tc.maxFileSize)
			lock, err := db.AcquireLinkRLock(tc.linkID)
			if err != nil {
				t.Fatal(err)
			}

			for _, f := range tc.files {
				fs.SetListLinkFilesResponse(f.Name, f.Files)
			}

			files, err := fs.ListLinkFiles(lock.ID())
			if err != nil {
				t.Fatal(err)
			}

			fileViews := Files(lock, files, tc.offset, tc.limit)

			if len(fileViews) != len(tc.res) {
				t.Fatalf("expected %d rendered files, got %d", len(tc.res), len(fileViews))
			}

			for i := range fileViews {
				if !reflect.DeepEqual(fileViews[i], tc.res[i]) {
					t.Fatalf(`incorrect file view at index %d:
expected: %v
got:      %v`, i, tc.res[i], fileViews[i])
				}
			}
		})
	}
}
