package flash

import "net/http"

type FlashKind int

const (
	SuccessFlash FlashKind = iota
	ErrorFlash
	FlashesCount
)

func (f FlashKind) String() string {
	switch f {
	case SuccessFlash:
		return "success_flash"
	case ErrorFlash:
		return "error_flash"
	default:
		return ""
	}
}

func AddFlash(w http.ResponseWriter, kind FlashKind, text string) {
	http.SetCookie(w, &http.Cookie{
		Name:   kind.String(),
		Path:   "/",
		Value:  text,
		MaxAge: 60,
	})
}

func GetFlashes(w http.ResponseWriter, r *http.Request) map[FlashKind]string {
	flashes := make(map[FlashKind]string)

	for kind := range FlashesCount {
		cookie, err := r.Cookie(kind.String())
		if err == nil {
			flashes[kind] = cookie.Value
		}

		http.SetCookie(w, &http.Cookie{
			Name:   kind.String(),
			Path:   "/",
			MaxAge: -1,
		})
	}

	return flashes
}
