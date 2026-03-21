package view

import (
	"reflect"
	"testing"
	"time"

	"github.com/foxpy/send-me-the-data/cmd/server/ifs"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs/mockfs"
	"github.com/foxpy/send-me-the-data/cmd/server/templates"
)

func TestFiles(t *testing.T) {
	for _, tc := range []struct {
		desc                string
		linkID              string
		renderDownloadLinks bool
		files               []linkFiles
		res                 []templates.FileView
	}{
		{
			desc:                "no files",
			linkID:              "abcd",
			renderDownloadLinks: false,
			files: []linkFiles{{
				name:  "abcd",
				files: []ifs.File{},
			}},
			res: []templates.FileView{},
		},
		{
			desc:                "one file",
			linkID:              "abcd",
			renderDownloadLinks: false,
			files: []linkFiles{{
				name: "abcd",
				files: []ifs.File{{
					Name:    "file 1",
					Size:    1024,
					ModTime: time.UnixMicro(0).UTC(),
				}},
			}},
			res: []templates.FileView{{
				Name:         "file 1",
				UploadedAt:   "Jan  1 00:00:00",
				Size:         "1.00 KiB",
				DownloadLink: "",
				DeleteLink:   "/link/abcd/file/file 1/delete",
			}},
		},
		{
			desc:                "render download link",
			linkID:              "abcd",
			renderDownloadLinks: true,
			files: []linkFiles{{
				name: "abcd",
				files: []ifs.File{{
					Name:    "file 1",
					Size:    1024,
					ModTime: time.UnixMicro(0).UTC(),
				}},
			}},
			res: []templates.FileView{{
				Name:         "file 1",
				UploadedAt:   "Jan  1 00:00:00",
				Size:         "1.00 KiB",
				DownloadLink: "/link/abcd/file/file 1",
				DeleteLink:   "/link/abcd/file/file 1/delete",
			}},
		},
		{
			desc:                "many files",
			linkID:              "abcd",
			renderDownloadLinks: false,
			files: []linkFiles{{
				name: "abcd",
				files: []ifs.File{
					{
						Name:    "file 1",
						Size:    1024,
						ModTime: time.UnixMicro(0).UTC(),
					},
					{
						Name:    "file 2",
						Size:    512,
						ModTime: time.UnixMicro(0).UTC(),
					},
					{
						Name:    "file 3",
						Size:    512,
						ModTime: time.UnixMicro(0).UTC(),
					},
				},
			}},
			res: []templates.FileView{
				{
					Name:         "file 1",
					UploadedAt:   "Jan  1 00:00:00",
					Size:         "1.00 KiB",
					DownloadLink: "",
					DeleteLink:   "/link/abcd/file/file 1/delete",
				},
				{
					Name:         "file 2",
					UploadedAt:   "Jan  1 00:00:00",
					Size:         "512 bytes",
					DownloadLink: "",
					DeleteLink:   "/link/abcd/file/file 2/delete",
				},
				{
					Name:         "file 3",
					UploadedAt:   "Jan  1 00:00:00",
					Size:         "512 bytes",
					DownloadLink: "",
					DeleteLink:   "/link/abcd/file/file 3/delete",
				},
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			fs := mockfs.NewMockFS()

			for _, f := range tc.files {
				fs.SetListLinkFilesResponse(f.name, f.files)
			}

			fileViews, err := Files(fs, tc.linkID, tc.renderDownloadLinks)
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
