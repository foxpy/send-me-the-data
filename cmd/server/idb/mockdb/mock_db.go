package mockdb

import "github.com/foxpy/send-me-the-data/cmd/server/idb"

type MockDB struct {
	allLinksResponse []idb.Link
}

var _ idb.Database = &MockDB{}

func NewMockDB() *MockDB {
	return &MockDB{}
}

func (d *MockDB) GetFileJournalEntry() (*idb.FileJournalEntry, error) {
	panic("TODO")
}

func (d *MockDB) DeleteFileJournalEntry(*idb.FileJournalEntry) error {
	panic("TODO")
}

func (d *MockDB) CreateFileJournalEntry(*idb.FileJournalEntry) error {
	panic("TODO")
}

func (d *MockDB) CreateLink(name, externalKey string, userDownloadable bool) error {
	panic("TODO")
}

func (d *MockDB) AcquireLinkRLock(externalKey string) (idb.LinkRLock, error) {
	panic("TODO")
}

func (d *MockDB) AcquireLinkWLock(externalKey string) (idb.LinkWLock, error) {
	panic("TODO")
}
