package mockdb

import (
	"fmt"

	"github.com/foxpy/send-me-the-data/src/idb"
)

type MockDB struct {
	listLinksResponses        map[struct{ offset, limit uint }][]idb.Link
	totalLinksResponse        *uint64
	expectedCreateLinkCalls   []CreateLinkCall
	acquireLinkRLockResponses map[string]idb.LinkRLock
}

type CreateLinkCall struct {
	link       link
	resultFunc func() error
}

var _ idb.Database = &MockDB{}

func NewMockDB() *MockDB {
	return &MockDB{
		listLinksResponses:        make(map[struct{ offset, limit uint }][]idb.Link),
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
	panic("not implemented")
}

func (d *MockDB) DeleteFileJournalEntry(*idb.FileJournalEntry) error {
	panic("not implemented")
}

func (d *MockDB) CreateFileJournalEntry(*idb.FileJournalEntry) error {
	panic("not implemented")
}

func (d *MockDB) AcquireLinkWLock(string) (idb.LinkWLock, error) {
	panic("not implemented")
}

func (d *MockDB) GetAdminPasswordHash(username string) ([]byte, error) {
	panic("not implemented")
}

func (d *MockDB) CreateSessionToken(token idb.SessionToken) error {
	panic("not implemented")
}

func (d *MockDB) GetSessionToken(token string) (*idb.SessionToken, error) {
	panic("not implemented")
}

func (d *MockDB) DeleteSessionToken(token string) error {
	panic("not implemented")
}

func (d *MockDB) DeleteOutdatedSessionTokens() error {
	panic("not implemented")
}

func (d *MockDB) CreateAdmin(username string, passwordHash []byte) error {
	panic("not implemented")
}

func (d *MockDB) DeleteAdmin(username string) error {
	panic("not implemented")
}
