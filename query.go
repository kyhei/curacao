package curacao

import (
	"errors"
)

// HTTPQuery URLに添付されているクエリ文字列について
type HTTPQuery map[string]string

// NewHTTPQuery convert HTTP query string to map[string]string
func NewHTTPQuery(q map[string][]string) HTTPQuery {
	query := make(HTTPQuery)

	for k, v := range q {
		query[k] = v[0]
	}

	return query
}

// Get get HTTPQuery value
func (p HTTPQuery) Get(name string) (string, error) {
	if val, ok := p[name]; ok {
		return val, nil
	}

	return "", errors.New("invalid named params")
}
