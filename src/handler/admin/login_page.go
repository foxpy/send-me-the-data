package admin

import (
	"net/http"

	"github.com/foxpy/send-me-the-data/src/flash"
	"github.com/foxpy/send-me-the-data/src/template"
)

func (s *AdminServer) loginPage(w http.ResponseWriter, r *http.Request) error {
	var params template.Params[struct{}]
	params.Title = "Send me the Data: Log in"

	flashes := flash.GetFlashes(w, r)
	params.SuccessFlash = flashes[flash.SuccessFlash]
	params.ErrorFlash = flashes[flash.ErrorFlash]

	return template.RenderAdminLogin(w, params)
}
