package web

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// A Handler is a type that handles an http request within our own little mini
// framework.
type Handler func(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params)

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct
type App struct {
	*httprouter.Router
	log *log.Logger
	mw  []Middleware
}

// New creates an App value that handle a set of routes for the application.
func New(log *log.Logger, mw ...Middleware) *App {
	return &App{
		Router: httprouter.New(),
		log:    log,
		mw:     mw,
	}
}

// Handle is our mechanism for mounting Handlers for a given HTTP verb and path
// pair, this makes for really easy, convenient routing.
func (a *App) Handle(verb, path string, handler Handler, mw ...Middleware) {

	// Wrap up the application-wide first, this will call the first function
	// of each middleware which will return a function of type Handler.``
	handler = wrapMiddleware(wrapMiddleware(handler, mw), a.mw)

	// The function to execute for each request.
	h := func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		// Setup CORS for development - Dont do this unless its a public endpoint obviously
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		// Call the wrapped handler functions.
		handler(a.log, w, r, params)
	}
	// Add this handler for the specified verb and route.
	a.Router.Handle(verb, path, h)

	//Handle preflight browser requests when we are POSTing data
	// f := func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// 	// Setup CORS for dev only or public endpoint
	// 	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	// 	return
	// }
	// a.Router.Handle("OPTIONS", path, f)
}
