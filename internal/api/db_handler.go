package api

import "database/sql"

// DBHandler is a parent struct that implements the HandlerFuncs so that they have access to the database handle.
type DBHandler struct {
	DB *sql.DB
}
