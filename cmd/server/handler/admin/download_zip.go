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

	// TODO: maybe I can achieve better speeds registering DEFLATE with lower compression level
	zw := zip.NewWriter(w)
	zw.AddFS(linkFS)
	err = zw.Close()
	if err != nil {
		return fmt.Errorf("failed to create ZIP archive for link %s: %w", lock.ExternalKey(), err)
	}

	return nil
}
