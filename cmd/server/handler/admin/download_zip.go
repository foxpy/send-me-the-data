package admin

import (
	"archive/zip"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/foxpy/send-me-the-data/cmd/server/handler"
)

// TODO: allow downloading uncompressed ZIP

func (s *AdminServer) downloadZIP(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	lock, err := s.db.AcquireLinkRLock(id)
	if errors.Is(err, sql.ErrNoRows) {
		return handler.Respond404(w)
	} else if err != nil {
		return fmt.Errorf("failed to acquire read lock for link %s: %w", id, err)
	}

	defer lock.Release()

	linkFS, err := s.fs.LinkFS(lock.ExternalKey())
	if err != nil {
		return fmt.Errorf("failed to obtain FS for link %s: %w", lock.ExternalKey(), err)
	}

	w.Header().Add("Content-Type", "application/zip")
	w.Header().Add(
		"Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s.zip"`, url.PathEscape(lock.Name())),
	)
	w.WriteHeader(http.StatusOK)

	zw := zip.NewWriter(w)
	err = zw.AddFS(linkFS)
	if err != nil {
		return fmt.Errorf("failed to create ZIP archive for link %s: %w", lock.ExternalKey(), err)
	}

	err = zw.Close()
	if err != nil {
		return fmt.Errorf("failed to finalize ZIP archive for link %s: %w", lock.ExternalKey(), err)
	}

	return nil
}
