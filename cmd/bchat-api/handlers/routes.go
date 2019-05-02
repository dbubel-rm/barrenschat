package handlers

import (
	"fmt"
	"log"
	"net/http"

	mid "github.com/dbubel/bchat/internal/middleware"
	"github.com/dbubel/bchat/internal/platform/db"
	"github.com/dbubel/bchat/internal/platform/web"
)

// API returns a handler for a set of routes.
func API(log *log.Logger, db *db.SQLite) http.Handler {

	// Application is created here along with adding any API wide middlewares
	var app = web.New(log, mid.RequestLogger, mid.Mid1, mid.Mid2)

	//
	// Each section of the API is established here
	//

	// Health
	check := Health{
		MasterDB: db.Database,
	}
	app.Router.NotFound = check.Health()

	app.Handle("POST", "/health", check.Health)

	// Traces
	uploads := Traces{
		MasterDB: db.Database,
	}

	app.Handle("GET", "/v1/getports", uploads.getPorts)

	return app
}

func helloworld2() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello World")
	})
}
