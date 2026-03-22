package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/foxpy/send-me-the-data/cmd/server/handler"
	"github.com/foxpy/send-me-the-data/cmd/server/idb"
	"github.com/foxpy/send-me-the-data/cmd/server/ifs"
)

func (s *AdminServer) downloadFile(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	name, err := handler.SanitizeFileName(r.PathValue("name"))
	if err != nil {
		return err
	}

	file, err := prepareDownloadFile(s.db, s.fs, id, name)
	if errors.Is(err, sql.ErrNoRows) {
		return handler.Respond404(w)
	} else if err != nil {
		return fmt.Errorf("failed to prepare file %s in link %s for serving: %w", name, id, err)
	}

	defer func() {
		_ = file.Close()
	}()

	http.ServeContent(w, r, name, time.Time{}, file)
	return nil
}

func prepareDownloadFile(db idb.Database, fs ifs.Filesystem, id, name string) (io.ReadSeekCloser, error) {
	lock, err := db.AcquireLinkRLock(id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("failed to acquire read lock for link %s: %w", id, err)
	}

	defer func() {
		_ = lock.Release()
	}()

	linkFS, err := fs.LinkFS(id)
	if err != nil {
		return nil, fmt.Errorf("failed to open filesystem to serve file: %w", err)
	}

	file, err := linkFS.Open(name)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s from link %s for reading: %w", name, id, err)
	}

	rsc, ok := file.(io.ReadSeekCloser)
	if !ok {
		return nil, errors.New("file is not io.ReadSeekCloser")
	}

	return rsc, err
}
