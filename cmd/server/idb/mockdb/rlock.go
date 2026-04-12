package mockdb

import (
	"time"

	"github.com/foxpy/send-me-the-data/cmd/server/idb"
)

type mockLinkRLock struct {
	name, externalKey string
	createdAt         time.Time
	userDowbnloadable bool
	maxFileSize       uint64
}

func (l *mockLinkRLock) UserDownloadable() bool {
	return l.userDowbnloadable
}

func (l *mockLinkRLock) Name() string {
	return l.name
}

func (l *mockLinkRLock) CreatedAt() time.Time {
	return l.createdAt
}

func (l *mockLinkRLock) ExternalKey() string {
	return l.externalKey
}

func (l *mockLinkRLock) MaxFileSize() uint64 {
	return l.maxFileSize
}

func (l *mockLinkRLock) Release() error {
	panic("TODO")
}

func (d *MockDB) AcquireLinkRLock(externalKey string) (idb.LinkRLock, error) {
	lock, ok := d.acquireLinkRLockResponses[externalKey]
	if !ok {
		panic("must mock AcquireLinkRLock() response")
	}

	return lock, nil
}

func (d *MockDB) SetAcquireLinkRLockResponse(
	externalKey string,
	name string,
	createdAt time.Time,
	userDownloadable bool,
	maxFileSize uint64,
) {
	d.acquireLinkRLockResponses[externalKey] = &mockLinkRLock{
		name, externalKey, createdAt, userDownloadable, maxFileSize,
	}
}
