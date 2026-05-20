package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/foxpy/send-me-the-data/src/idb"
	"github.com/foxpy/send-me-the-data/src/ifs"
)

var errNotReadSeekCloser = errors.New("file is not io.ReadSeekCloser")

func HandleDownloadFile(
	w http.ResponseWriter,
	r *http.Request,
	db idb.Database,
	fs ifs.Filesystem,
	checkUserDownloadable bool,
) error {
	id := r.PathValue("id")
	name, err := SanitizeFileName(r.PathValue("name"))
	if err != nil {
		return err
	}

	lock, err := db.AcquireLinkRLock(id)
	if errors.Is(err, sql.ErrNoRows) {
		return RespondError(w, http.StatusNotFound)
	} else if err != nil {
		return fmt.Errorf("failed to acquire read lock for link %s: %w", id, err)
	}

	defer lock.Release()

	if checkUserDownloadable && !lock.UserDownloadable() {
		return RespondError(w, http.StatusForbidden)
	}

	linkFS, err := fs.LinkFS(id)
	if err != nil {
		return fmt.Errorf("failed to open filesystem to serve file: %w", err)
	}

	file, err := linkFS.Open(name)
	if err != nil {
		return fmt.Errorf("failed to open file %s from link %s for reading: %w", name, id, err)
	}

	defer func() {
		_ = file.Close()
	}()

	lock.Release()

	rsc, ok := file.(io.ReadSeekCloser)
	if !ok {
		return errNotReadSeekCloser
	}

	http.ServeContent(w, r, name, time.Time{}, rsc)
	return nil
}
