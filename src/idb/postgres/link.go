package postgres

import "time"

type link struct {
	name, id                        string
	createdAt                       time.Time
	userDownloadable, uploadEnabled bool
	maxFileSize                     uint64
}

func (l *link) Name() string {
	return l.name
}

func (l *link) ID() string {
	return l.id
}

func (l *link) CreatedAt() time.Time {
	return l.createdAt
}

func (l *link) UserDownloadable() bool {
	return l.userDownloadable
}

func (l *link) UploadEnabled() bool {
	return l.uploadEnabled
}

func (l *link) MaxFileSize() uint64 {
	return l.maxFileSize
}
