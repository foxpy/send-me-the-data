package idb

import (
	"time"
)

type Database interface {
	GetFileJournalEntry() (*FileJournalEntry, error)
	DeleteFileJournalEntry(*FileJournalEntry) error
	CreateFileJournalEntry(*FileJournalEntry) error
	TotalLinks() (uint64, error)
	ListLinks(offset, limit uint) ([]Link, error)
	CreateLink(name, id string, userDownloadable, uploadEnabled bool, maxFileSize uint64) error
	AcquireLinkRLock(id string) (LinkRLock, error)
	AcquireLinkWLock(id string) (LinkWLock, error)
	GetAdminPasswordHash(username string) ([]byte, error)
	CreateSessionToken(token SessionToken) error
	GetSessionToken(token string) (*SessionToken, error)
	DeleteSessionToken(token string) error
	DeleteOutdatedSessionTokens() error
	CreateAdmin(username string, passwordHash []byte) error
	DeleteAdmin(username string) error
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
	Release()
}

type LinkWLock interface {
	Link
	Update(name string, userDownloadable, uploadEnabled bool, maxFileSize uint64) error
	Delete() error
	Commit() error
	Rollback()
}

type SessionToken struct {
	Token     string
	Username  string
	ExpiresAt time.Time
}
