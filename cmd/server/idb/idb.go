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
	CreateLink(name, externalKey string, userDownloadable bool) error
	AcquireLinkRLock(externalKey string) (LinkRLock, error)
	AcquireLinkWLock(externalKey string) (LinkWLock, error)
	// TODO: this function doesn't really belong here
	GenerateRandomExternalKey() string
}

type FileJournalEntry struct {
	LinkExternalKey string
	FileName        string
}

type Link struct {
	Name, ExternalKey string
	CreatedAt         time.Time
	UserDownloadable  bool
}

type LinkRLock interface {
	UserDownloadable() bool
	Name() string
	CreatedAt() time.Time
	ExternalKey() string
	Release() error
}

type LinkWLock interface {
	UpdateLink(name string, userDownloadable bool) error
	DeleteLink() error
	Release() error
}
