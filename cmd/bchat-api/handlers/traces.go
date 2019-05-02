package handlers

import (
	"log"
	"net/http"

	"github.com/dbubel/bchat/internal/platform/web"
	"github.com/dbubel/bchat/internal/traces"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

type Traces struct {
	MasterDB *sqlx.DB
}

// H is a shortcut for map[string]interface{}
type H map[string]interface{}

func (c *Traces) getPorts(log *log.Logger, w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	ports, err := traces.GetPorts(log, c.MasterDB, params.ByName("hostId"))
	if err != nil {
		web.RespondError(log, w, err, http.StatusInternalServerError)
		return
	}

	web.Respond(log, w, ports, http.StatusOK)
}
