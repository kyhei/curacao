package curacao

import (
	"net/http"
	"regexp"
	"strings"
)

// HTTPHandler the type
type HTTPHandler interface{}

// MiddlewareResponse response type of ExecuteMiddleware
type MiddlewareResponse struct {
	OK         bool
	Code       int
	Propagates interface{}
}

func newMiddlewareResponse() *MiddlewareResponse {
	mr := new(MiddlewareResponse)
	mr.OK = true
	mr.Code = 200
	mr.Propagates = nil
	return mr
}

// Dispatcher これは所謂Dispatcher
type Dispatcher struct {
	Routes     []Route
	Middleware []interface{}
}

// NewDispatcher 新規Dispatcherを作成
func newDispatcher() *Dispatcher {
	dispatcher := new(Dispatcher)
	return dispatcher
}

// Route Dispatcherを管理するもの
type Route struct {
	method  string
	reg     *regexp.Regexp
	handler HTTPHandler
}

func newRoute(method string, path string, handler HTTPHandler) (*Route, bool) {
	route := new(Route)
	route.method = strings.ToUpper(method)
	route.handler = handler
	route.reg = regexp.MustCompile(path)

	return route, ValidateHandler(handler)
}

func (d *Dispatcher) register(method string, path string, handler HTTPHandler) {
	convertedPath := convertURLString(path)
	route, ok := newRoute(method, convertedPath, handler)
	if ok == false {
		panic("the type of handler is must be func")
	}

	routes := append(d.Routes, *route)
	d.Routes = routes
}

func (d *Dispatcher) use(middleware interface{}) {
	ok := ValidateMiddleware(middleware)
	if ok == false {
		panic("middleware is invalid")
	}
	d.Middleware = append(d.Middleware, middleware)
}

func (d *Dispatcher) dispatch(r *http.Request) (
	HTTPHandler,
	HTTPParams,
	HTTPQuery,
) {

	code := http.StatusNotFound

	for _, route := range d.Routes {
		if route.reg.MatchString(r.URL.Path) {

			if route.method == r.Method {
				return route.handler,
					NewHTTPParams(r.URL.Path, route.reg),
					NewHTTPQuery(r.URL.Query())
			}

			code = http.StatusMethodNotAllowed

		}
	}

	return func() int { return code }, nil, nil
}

// Router httpリクエストをルーティングします。
func router(dispatcher Dispatcher) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		f, p, q := dispatcher.dispatch(r)

		mrs := []interface{}{}
		for _, m := range dispatcher.Middleware {
			mr := ExecuteMiddleware(m, r)
			if mr.OK == false {
				Render(w, int(mr.Code), nil)
				return
			}

			mrs = append(mrs, mr.Propagates)
		}

		args := []interface{}{r, p, q}
		args = append(args, mrs...)
		code, resp := ExecuteHandler(f, args)
		Render(w, int(code), resp)
	}
}

func convertURLString(path string) string {

	cleanPath := path
	if strings.HasSuffix(path, "/") {
		cleanPath = path[:len(path)-1]
	}

	params := make([]string, 0)
	for index, s := range strings.Split(cleanPath, "/") {
		if index == 0 || s == "" {
			continue
		}

		if strings.HasPrefix(s, ":") {
			params = append(params, s[1:])
		}

	}

	cleanPath = strings.Replace(cleanPath, "/", `\/`, -1)

	for _, param := range params {
		cleanPath = strings.Replace(cleanPath, ":"+param, `(?P<`+param+`>[^\/]*)`, 1)
	}

	res := `^` + cleanPath + `\/?$`

	return res
}
