package curacao

import (
	"net/http"
)

// Render redering http response
func Render(
	w http.ResponseWriter,
	code int,
	body []byte,
) {
	w.WriteHeader(code)

	if body != nil {
		w.Write(body)
	}
}
