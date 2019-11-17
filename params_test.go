package curacao

import (
	"regexp"
	"testing"
)

func TestParams(t *testing.T) {
	path := "/params/any/this-is-message/"
	reg := regexp.MustCompile(`\/params\/any\/(?P<message>[^\/]*)\/?$`)

	httpParams := NewHTTPParams(path, reg)

	_, err := httpParams.Get("message")
	if err != nil {
		t.Errorf("failed parse path string. failed path is %v", path)
	}

	_, err = httpParams.Get("msg")
	if err == nil {
		t.Errorf("failed parse path string. failed path is %v", path)
	}

	path = "/hoge/fuga/moja/"
	reg = regexp.MustCompile(`\/(?P<p1>[^\/]*)\/(?P<p2>[^\/]*)\/(?P<p3>[^\/]*)\/?$`)
	httpParams = NewHTTPParams(path, reg)

	p1, err := httpParams.Get("p1")
	if p1 != "hoge" {
		t.Errorf("p1 params does not match 'hoge', actual: %v", p1)
	}

	p2, err := httpParams.Get("p2")
	if p2 != "fuga" {
		t.Errorf("p2 params does not match 'fuga', actual: %v", p2)
	}

	p3, err := httpParams.Get("p3")
	if p3 != "moja" {
		t.Errorf("p3 params does not match 'moja', actual: %v", p3)
	}

	_, err = httpParams.Get("p4")
	if err == nil {
		t.Error("p4 params does not exsist. But found it")
	}

}
