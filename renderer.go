package curacao

import (
	"net/http"
)

func render(
	w http.ResponseWriter,
	code int,
	body []byte,
) {
	w.WriteHeader(code)

	if body != nil {
		w.Write(body)
	}
}
