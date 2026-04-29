package mockdb

import (
	"time"

	"github.com/foxpy/send-me-the-data/cmd/server/idb"
)

type link struct {
	name, id                        string
	createdAt                       time.Time
	userDownloadable, uploadEnabled bool
	maxFileSize                     uint64
}

func (l *link) UserDownloadable() bool {
	return l.userDownloadable
}

func (l *link) UploadEnabled() bool {
	return l.uploadEnabled
}

func (l *link) Name() string {
	return l.name
}

func (l *link) CreatedAt() time.Time {
	return l.createdAt
}

func (l *link) ID() string {
	return l.id
}

func (l *link) MaxFileSize() uint64 {
	return l.maxFileSize
}

func NewLink(
	id string,
	name string,
	createdAt time.Time,
	userDownloadable bool,
	uploadEnabled bool,
	maxFileSize uint64,
) idb.Link {
	return &link{name, id, createdAt, userDownloadable, uploadEnabled, maxFileSize}
}
