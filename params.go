package curacao

import (
	"errors"
	"regexp"
)

// HTTPParams HTTP paramaters
type HTTPParams map[string]string

// NewHTTPParams 新しいNewHTTPParamsを作成します。
func NewHTTPParams(path string, reg *regexp.Regexp) HTTPParams {
	params := make(HTTPParams)

	match := reg.FindStringSubmatch(path)

	for index, name := range reg.SubexpNames() {
		if index != 0 {
			params[name] = match[index]
		}
	}

	return params
}

// Get get HTTPParams value
func (p HTTPParams) Get(name string) (string, error) {
	if val, ok := p[name]; ok {
		return val, nil
	}

	return "", errors.New("invalid named params")
}
