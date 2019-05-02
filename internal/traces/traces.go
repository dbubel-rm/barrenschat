package traces

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func GetPorts(log *log.Logger, dbConn *sqlx.DB, hostID string) ([]Port, error) {
	// TODO: these results could be cached.
	var queryReturn []Port
	err := dbConn.Select(&queryReturn, "SELECT * FROM ports where host_id = ?", hostID)

	if err != nil {
		log.Println("Error querying for ports:", err.Error())
		return queryReturn, err
	}

	return queryReturn, nil
}
