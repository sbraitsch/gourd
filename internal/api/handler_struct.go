package api

import "database/sql"

// HandlerStruct is a parent struct that implements the HandlerFuncs so that they have access to the database handle.
type HandlerStruct struct {
	DB *sql.DB
}
