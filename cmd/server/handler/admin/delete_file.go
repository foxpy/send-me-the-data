package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/foxpy/send-me-the-data/cmd/server/handler"
)

func (s *AdminServer) deleteFile(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	lock, err := s.db.AcquireLinkRLock(id)
	if errors.Is(err, sql.ErrNoRows) {
		return handler.Respond404(w)
	} else if err != nil {
		return fmt.Errorf("failed to acquire read lock for link %s: %w", id, err)
	}

	defer func() {
		_ = lock.Release()
	}()

	name := r.PathValue("name")
	err = s.fs.RemoveLinkFile(id, name)
	if err != nil {
		return fmt.Errorf("failed to remove file %s from link %s: %w", name, id, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/link/%s", id), http.StatusSeeOther)
	return nil
}
