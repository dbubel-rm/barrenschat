package handlers

import (
	"log"
	"net/http"

	"github.com/dbubel/bchat/internal/platform/web"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

// Health provides support for orchestration health checks.
type Health struct {
	MasterDB *sqlx.DB
}

// Health validates the service is healthy and ready to accept requests.
func (c *Health) Health(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	status := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	err := c.MasterDB.Ping()
	if err != nil {
		web.RespondError(log, w, err, http.StatusInternalServerError)
	}

	web.Respond(log, w, status, http.StatusOK)

}
