package idb

import (
	"time"
)

type Database interface {
	GetFileJournalEntry() (*FileJournalEntry, error)
	DeleteFileJournalEntry(*FileJournalEntry) error
	CreateFileJournalEntry(*FileJournalEntry) error
	// FIXME: do not read all links from database, use pagination instead
	AllLinks() ([]Link, error)
	CreateLink(name, id string, userDownloadable, uploadEnabled bool, maxFileSize uint64) error
	AcquireLinkRLock(id string) (LinkRLock, error)
	AcquireLinkWLock(id string) (LinkWLock, error)
	// TODO: this function doesn't really belong here
	GenerateRandomPublicID() string
}

type FileJournalEntry struct {
	LinkPublicID string
	FileName     string
}

type Link interface {
	Name() string
	ID() string
	CreatedAt() time.Time
	UserDownloadable() bool
	UploadEnabled() bool
	MaxFileSize() uint64
}

type LinkRLock interface {
	Link
	Release() error
}

type LinkWLock interface {
	Link
	Update(name string, userDownloadable, uploadEnabled bool, maxFileSize uint64) error
	Delete() error
	Commit() error
	Rollback() error
}
