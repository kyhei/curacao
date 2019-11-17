package curacao

import (
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
	"unsafe"
)

func Test_ConvertURLString(t *testing.T) {
	rawPath := "/"
	expectedPath := `^\/?$`
	convertedPath := convertURLString(rawPath)

	if convertedPath != expectedPath {
		t.Errorf("converted Path failed. expected: %v got: %v", expectedPath, convertedPath)
	}

	rawPath = "/hoge/:message_1/fuga/:message_2"
	expectedPath = `^\/hoge\/(?P<message_1>[^\/]*)\/fuga\/(?P<message_2>[^\/]*)\/?$`
	convertedPath = convertURLString(rawPath)

	if convertedPath != expectedPath {
		t.Errorf("converted Path failed. expected: %v got: %v", expectedPath, convertedPath)
	}

}

func Test_Router(t *testing.T) {
	c := NewApp("localhost", "8080")
	c.Get("/", func() string { return "hello curacao!" })
	c.Post("/", func() string { return "POST /" })
	c.Get("/params/:message", func(p HTTPParams) string {
		param, err := p.Get("message")
		if err != nil {
			t.Error("curacao internal server error")
		}
		return param
	})
	c.Get("/pages/:id/show", func() string { return "GET /pages/:id/show" })
	c.Get("/many/:hoge/:fuga", func() string { return "GET /many/:hoge/:fuga" })
	go c.Start()

	time.Sleep(1 * time.Second)

	client := &http.Client{}

	resp, err := client.Get("http://localhost:8080/")
	if err != nil {
		t.Error("http request is failed")
	}

	bs := assertToString(resp.Body)
	resp.Body.Close()

	if bs != "hello curacao!" {
		t.Errorf("expect response is 'hello curacao!', got: %v", bs)
	}

	resp, err = client.Post("http://localhost:8080/", "application/json", nil)
	if err != nil {
		t.Error("http request is failed")
	}

	bs = assertToString(resp.Body)
	resp.Body.Close()

	if bs != "POST /" {
		t.Errorf("expect response is 'POST /', got: %v", bs)
	}

	resp, err = client.Get("http://localhost:8080/hoge")
	if err != nil {
		t.Error("http request is failed")
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code is %v, got: %v", http.StatusNotFound, resp.StatusCode)
	}

	resp, err = client.Get("http://localhost:8080/params/curacao")
	if err != nil {
		t.Error("http request is failed")
	}

	bs = assertToString(resp.Body)
	resp.Body.Close()

	if bs != "curacao" {
		t.Errorf("expected response is curacao, got: %v", bs)
	}

	resp, err = client.Get("http://localhost:8080/pages/1123/show?preview=true")
	if err != nil {
		t.Error("http request is failed")
	}

	bs = assertToString(resp.Body)
	resp.Body.Close()

	if bs != "GET /pages/:id/show" {
		t.Errorf("expected response is 'GET /pages/:id/show', got: %v", bs)
	}

	resp, err = client.Get("http://localhost:8080/many/moja/curacao")
	if err != nil {
		t.Error("http request is failed")
	}

	bs = assertToString(resp.Body)
	resp.Body.Close()

	if bs != "GET /many/:hoge/:fuga" {
		t.Errorf("expected response is 'GET /many/:hoge/:fuga', got: %v", bs)
	}

}

func assertToString(i io.Reader) string {
	b, err := ioutil.ReadAll(i)
	if err != nil {
		panic("failed to assert")
	}
	bs := *(*string)(unsafe.Pointer(&b))
	return bs
}
