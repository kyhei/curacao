package curacao

import (
	"net/http"
)

func render(
	w http.ResponseWriter,
	header map[string]string,
	code int,
	body []byte,
) {

	for name, value := range header {
		w.Header().Set(name, value)
	}

	w.WriteHeader(code)

	if body != nil {
		w.Write(body)
	}
}
