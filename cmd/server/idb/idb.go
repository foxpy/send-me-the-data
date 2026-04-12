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
	CreateLink(name, externalKey string, userDownloadable bool, maxFileSize uint64) error
	AcquireLinkRLock(externalKey string) (LinkRLock, error)
	AcquireLinkWLock(externalKey string) (LinkWLock, error)
	// TODO: this function doesn't really belong here
	GenerateRandomExternalKey() string
}

type FileJournalEntry struct {
	LinkExternalKey string
	FileName        string
}

// TODO: I think it is better to make Link an interface, too,
//       as it will allow me to erase a lot of repeating boilerplate

type Link struct {
	Name, ExternalKey string
	CreatedAt         time.Time
	UserDownloadable  bool
	MaxFileSize       uint64
}

type LinkRLock interface {
	UserDownloadable() bool
	Name() string
	CreatedAt() time.Time
	ExternalKey() string
	MaxFileSize() uint64
	Release() error
}

type LinkWLock interface {
	UpdateLink(name string, userDownloadable bool, maxFileSize uint64) error
	DeleteLink() error
	Release() error
}
