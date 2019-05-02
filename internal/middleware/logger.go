package mid

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dbubel/bchat/internal/platform/web"
	"github.com/julienschmidt/httprouter"
)

// RequestLogger writes some information about the request to the logs in
func RequestLogger(before web.Handler) web.Handler {
	// Wrap this handler around the next one provided.
	return func(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		fmt.Println("I'm mid0")
		before(log, w, r, params)
		log.Printf("%s -> %d -> %s -> %s", r.Method, r.ContentLength, r.URL.Path, r.RemoteAddr)
	}
}

// RequestLogger writes some information about the request to the logs in
func Mid1(before web.Handler) web.Handler {
	// Wrap this handler around the next one provided.
	return func(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		fmt.Println("I'm mid1")
		before(log, w, r, params)
	}
}

// RequestLogger writes some information about the request to the logs in
func Mid2(before web.Handler) web.Handler {
	// Wrap this handler around the next one provided.
	return func(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		fmt.Println("I'm mid2")
		before(log, w, r, params)
	}
}
