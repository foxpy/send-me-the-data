package user

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/foxpy/send-me-the-data/cmd/server/handler"
	"github.com/foxpy/send-me-the-data/cmd/server/template"
)

func (s *UserServer) viewLinkPage(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	title, files, err := handler.PrepareFilesView(s.db, s.fs, id, false)
	if errors.Is(err, sql.ErrNoRows) {
		return handler.Respond404(w)
	} else if err != nil {
		return err
	}

	var params template.Params[template.UserViewLinkParams]
	params.Title = title
	params.Data.Files = files

	_, err = r.Cookie("success_flash")
	if err == nil {
		params.SuccessFlash = "File uploaded successfully"
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "success_flash",
		MaxAge: -1,
	})

	return template.RenderUserViewLink(w, params)
}
