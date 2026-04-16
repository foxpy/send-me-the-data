package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/foxpy/send-me-the-data/cmd/server/handler"
)

func (s *AdminServer) deleteLink(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	// FIXME: since we do not attempt to terminate all active uploads,
	//        acquiring this lock might take possibly many hours of time.
	// There is an idea for another fix: we could download files from users
	// to a temporary designated directory and only take a lock when moving
	// the file from this directory to the target link directory. This way,
	// we can delete a link almost instantly, and when the upload endpoint
	// finishes, it will fail to move the file because it will fail to acquire
	// a lock on an unexisting link, therefore forced to delete a file it just
	// finished downloading. Still not perfect, but much better.
	lock, err := s.db.AcquireLinkWLock(id)
	if errors.Is(err, sql.ErrNoRows) {
		return handler.RespondError(w, http.StatusNotFound)
	} else if err != nil {
		return fmt.Errorf("failed to acquire write lock for link %s: %w", id, err)
	}

	defer func() {
		_ = lock.Commit()
	}()

	err = s.fs.RemoveLinkFiles(id)
	if err != nil {
		return err
	}

	err = lock.Delete()
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
