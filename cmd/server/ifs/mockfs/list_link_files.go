package mockfs

import "github.com/foxpy/send-me-the-data/cmd/server/ifs"

func (f *MockFS) ListLinkFiles(linkID string) ([]ifs.File, error) {
	resp, ok := f.listLinkFilesResponses[linkID]
	if !ok {
		panic("must mock ListLinkFiles() response")
	}

	return resp, nil
}

func (f *MockFS) SetListLinkFilesResponse(linkID string, response []ifs.File) {
	f.listLinkFilesResponses[linkID] = response
}
