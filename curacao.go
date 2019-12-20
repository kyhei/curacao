package curacao

import (
	"log"
	"net/http"
)

// App the tyoe of App
type App struct {
	host       string
	port       string
	dispatcher *dispatcher
}

// NewApp create New App
func NewApp(host string, port string) *App {
	app := new(App)
	app.host = host
	app.port = port
	app.dispatcher = newDispatcher()
	return app
}

// Register register new route
func (app *App) Register(method string, path string, handler HTTPHandler) {
	app.dispatcher.register(method, path, handler)
}

// RegisterWithPreflight register HTTP handler with OPTIONS method
func (app *App) RegisterWithPreflight(method string, path string, handler HTTPHandler) {
	app.Register(method, path, handler)
	app.Register("OPTIONS", path, func() int { return http.StatusOK })
}

// Use register middleware
func (app *App) Use(middleware interface{}) {
	app.dispatcher.use(middleware)
}

// Header set HTTP Response Header
func (app *App) Header(name string, value string) {
	app.dispatcher.setHeader(name, value)
}

// Get short hand of Register
func (app *App) Get(path string, handler HTTPHandler) {
	app.dispatcher.register("GET", path, handler)
}

// GetWithPreflight register HTTP handler with OPTIONS method
func (app *App) GetWithPreflight(path string, handler HTTPHandler) {
	app.Get(path, handler)
	app.Register("OPTIONS", path, func() int { return http.StatusOK })
}

// Post short hand of Register
func (app *App) Post(path string, handler HTTPHandler) {
	app.dispatcher.register("POST", path, handler)
}

// PostWithPreflight register HTTP handler with OPTIONS method
func (app *App) PostWithPreflight(path string, handler HTTPHandler) {
	app.Post(path, handler)
	app.Register("OPTIONS", path, func() int { return http.StatusOK })
}

// Start start HTTP server
func (app *App) Start() {
	http.HandleFunc("/", router(*app.dispatcher))
	log.Printf("starting server in http://%s:%s", app.host, app.port)
	log.Fatal(http.ListenAndServe(app.host+":"+app.port, nil))
}
