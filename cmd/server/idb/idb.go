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
	CreateLink(name, externalKey string, userDownloadable, uploadEnabled bool, maxFileSize uint64) error
	AcquireLinkRLock(externalKey string) (LinkRLock, error)
	AcquireLinkWLock(externalKey string) (LinkWLock, error)
	// TODO: this function doesn't really belong here
	GenerateRandomExternalKey() string
}

type FileJournalEntry struct {
	LinkExternalKey string
	FileName        string
}

type Link interface {
	Name() string
	ExternalKey() string
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
