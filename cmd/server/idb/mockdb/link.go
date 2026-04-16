package mockdb

import (
	"time"

	"github.com/foxpy/send-me-the-data/cmd/server/idb"
)

type link struct {
	name, externalKey               string
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

func (l *link) ExternalKey() string {
	return l.externalKey
}

func (l *link) MaxFileSize() uint64 {
	return l.maxFileSize
}

func NewLink(externalKey string,
	name string,
	createdAt time.Time,
	userDownloadable bool,
	uploadEnabled bool,
	maxFileSize uint64,
) idb.Link {
	return &link{name, externalKey, createdAt, userDownloadable, uploadEnabled, maxFileSize}
}
