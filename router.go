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

type dispatcher struct {
	Routes     []route
	Middleware []interface{}
	Header     map[string]string
}

func newDispatcher() *dispatcher {
	dispatcher := new(dispatcher)
	dispatcher.Header = make(map[string]string)
	return dispatcher
}

type route struct {
	method  string
	reg     *regexp.Regexp
	handler HTTPHandler
}

func newRoute(method string, path string, handler HTTPHandler) (*route, bool) {
	route := new(route)
	route.method = strings.ToUpper(method)
	route.handler = handler
	route.reg = regexp.MustCompile(path)

	return route, validateHandler(handler)
}

func (d *dispatcher) register(method string, path string, handler HTTPHandler) {
	convertedPath := convertURLString(path)
	route, ok := newRoute(method, convertedPath, handler)
	if ok == false {
		panic("the type of handler is must be func")
	}

	routes := append(d.Routes, *route)
	d.Routes = routes
}

func (d *dispatcher) use(middleware interface{}) {
	ok := validateMiddleware(middleware)
	if ok == false {
		panic("middleware is invalid")
	}
	d.Middleware = append(d.Middleware, middleware)
}

func (d *dispatcher) setHeader(name string, value string) {
	d.Header[name] = value
}

func (d *dispatcher) dispatch(r *http.Request) (
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

func router(dispatcher dispatcher) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		f, p, q := dispatcher.dispatch(r)

		mrs := []interface{}{}
		for _, m := range dispatcher.Middleware {
			mr := executeMiddleware(m, w, r)
			if mr.OK == false {
				render(w, dispatcher.Header, int(mr.Code), nil)
				return
			}

			mrs = append(mrs, mr.Propagates)
		}

		args := []interface{}{r, p, q}
		args = append(args, mrs...)
		code, resp := executeHandler(f, args)
		render(w, dispatcher.Header, int(code), resp)
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
