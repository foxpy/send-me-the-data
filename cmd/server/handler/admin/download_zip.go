package admin

import (
	"archive/zip"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"

	"github.com/foxpy/send-me-the-data/cmd/server/handler"
)

func (s *AdminServer) downloadZIP(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	lock, err := s.db.AcquireLinkRLock(id)
	if errors.Is(err, sql.ErrNoRows) {
		return handler.RespondError(w, http.StatusNotFound)
	} else if err != nil {
		return fmt.Errorf("failed to acquire read lock for link %s: %w", id, err)
	}

	defer lock.Release()

	var method uint16
	switch r.URL.Query().Get("method") {
	case "":
		fallthrough
	case "deflate":
		method = zip.Deflate
	case "store":
		method = zip.Store
	default:
		w.WriteHeader(http.StatusBadRequest)
		return nil
	}

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
	err = zipAddFS(zw, linkFS, method)
	if err != nil {
		return fmt.Errorf("failed to create ZIP archive for link %s: %w", lock.ExternalKey(), err)
	}

	err = zw.Close()
	if err != nil {
		return fmt.Errorf("failed to finalize ZIP archive for link %s: %w", lock.ExternalKey(), err)
	}

	return nil
}

func zipAddFS(w *zip.Writer, fsys fs.FS, method uint16) error {
	return fs.WalkDir(fsys, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if name == "." {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		if !d.IsDir() && !info.Mode().IsRegular() {
			return errors.New("zip: cannot add non-regular file")
		}
		h, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		h.Name = name
		if d.IsDir() {
			h.Name += "/"
		}
		h.Method = method
		fw, err := w.CreateHeader(h)
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		f, err := fsys.Open(name)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(fw, f)
		return err
	})
}
