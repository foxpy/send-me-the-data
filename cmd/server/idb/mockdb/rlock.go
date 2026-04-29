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

func (d *MockDB) AcquireLinkRLock(id string) (idb.LinkRLock, error) {
	lock, ok := d.acquireLinkRLockResponses[id]
	if !ok {
		panic("must mock AcquireLinkRLock() response")
	}

	return lock, nil
}

func (d *MockDB) SetAcquireLinkRLockResponse(
	id string,
	name string,
	createdAt time.Time,
	userDownloadable bool,
	uploadEnabled bool,
	maxFileSize uint64,
) {
	l := link{name, id, createdAt, userDownloadable, uploadEnabled, maxFileSize}
	d.acquireLinkRLockResponses[id] = &mockLinkRLock{l}
}
