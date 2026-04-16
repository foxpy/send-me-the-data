package mockdb

import (
	"time"

	"github.com/foxpy/send-me-the-data/cmd/server/idb"
)

type mockLinkRLock struct {
	link
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
	uploadEnabled bool,
	maxFileSize uint64,
) {
	l := link{name, externalKey, createdAt, userDownloadable, uploadEnabled, maxFileSize}
	d.acquireLinkRLockResponses[externalKey] = &mockLinkRLock{l}
}
