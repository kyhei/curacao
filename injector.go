package curacao

import (
	"net/http"
	"reflect"
	"unsafe"
)

// validate middleware before curacao start
func validateMiddleware(middleware interface{}) bool {
	fnt := reflect.TypeOf(middleware)
	if fnt.Kind() != reflect.Func {
		println("the type of middleware is must be func")
		return false
	}

	return true
}

// run curacao middleware
func executeMiddleware(middleware interface{}, args ...interface{}) *MiddlewareResponse {
	fnt := reflect.TypeOf(middleware)

	argsNum := fnt.NumIn()
	middlewareArgs := make([]reflect.Value, argsNum)

	writerType := reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()

	for i := 0; i < argsNum; i++ {
		argType := fnt.In(i)
		for _, arg := range args {
			if argType == reflect.TypeOf(arg) ||
				(reflect.TypeOf(arg).Implements(writerType) && argType.Implements(writerType)) {
				middlewareArgs[i] = reflect.ValueOf(arg)
				break
			}
		}
	}

	fnv := reflect.ValueOf(middleware)
	result := fnv.Call(middlewareArgs)
	response := newMiddlewareResponse()

	for i := 0; i < fnt.NumOut(); i++ {
		if r, ok := result[i].Interface().(int); ok {
			response.Code = r
			continue
		}

		if r, ok := result[i].Interface().(bool); ok {
			response.OK = r
			continue
		}

		if r, ok := result[i].Interface().(interface{}); ok {
			response.Propagates = r
			continue
		}
	}

	return response
}

// validate request handler before curacao start
func validateHandler(handler interface{}) bool {
	fnt := reflect.TypeOf(handler)
	if fnt.Kind() != reflect.Func {
		return false
	}

	return true
}

// run http handler function
func executeHandler(handler interface{}, args []interface{}) (int, []byte) {
	fnt := reflect.TypeOf(handler)

	argsNum := fnt.NumIn()
	handlerArgs := make([]reflect.Value, argsNum)

	for i := 0; i < argsNum; i++ {
		argType := fnt.In(i)
		for _, arg := range args {
			if argType == reflect.TypeOf(arg) {
				handlerArgs[i] = reflect.ValueOf(arg)
				break
			}
		}
	}

	fnv := reflect.ValueOf(handler)
	result := fnv.Call(handlerArgs)

	var body []byte
	code := 200

	for i := 0; i < len(result); i++ {
		if r, ok := result[i].Interface().(int); ok {
			code = r
			continue
		}

		if r, ok := result[i].Interface().(string); ok {
			body = *(*[]byte)(unsafe.Pointer(&r))
			continue
		}

		if r, ok := result[i].Interface().([]byte); ok {
			body = r
			continue
		}

	}

	return code, body

}
