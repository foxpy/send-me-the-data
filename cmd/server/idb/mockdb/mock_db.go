package mockdb

import (
	"fmt"

	"github.com/foxpy/send-me-the-data/cmd/server/idb"
)

type MockDB struct {
	allLinksResponse           []idb.Link
	randomExternalKeyResponses []string
	expectedCreateLinkCalls    []CreateLinkCall
	acquireLinkRLockResponses  map[string]idb.LinkRLock
}

type CreateLinkCall struct {
	link       idb.Link
	resultFunc func() error
}

var _ idb.Database = &MockDB{}

func NewMockDB() *MockDB {
	return &MockDB{
		acquireLinkRLockResponses: make(map[string]idb.LinkRLock),
	}
}

func (d *MockDB) CheckAllExpects() {
	if len(d.expectedCreateLinkCalls) > 0 {
		firstExpectedCall := d.expectedCreateLinkCalls[0]
		panic(fmt.Sprintf("expected CreateLink(%v) call, which never happened", firstExpectedCall))
	}
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

func (d *MockDB) AcquireLinkWLock(externalKey string) (idb.LinkWLock, error) {
	panic("TODO")
}
