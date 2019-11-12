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

// Use register middleware
func (app *App) Use(middleware interface{}) {
	app.dispatcher.use(middleware)
}

// Get short hand of Register
func (app *App) Get(path string, handler HTTPHandler) {
	app.dispatcher.register("GET", path, handler)
}

// Post short hand of Register
func (app *App) Post(path string, handler HTTPHandler) {
	app.dispatcher.register("POST", path, handler)
}

// Start start HTTP server
func (app *App) Start() {
	http.HandleFunc("/", router(*app.dispatcher))
	log.Printf("starting server in http://%s:%s", app.host, app.port)
	log.Fatal(http.ListenAndServe(app.host+":"+app.port, nil))
}
