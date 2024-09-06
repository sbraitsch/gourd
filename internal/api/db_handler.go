package api

import "database/sql"

type DBHandler struct {
	DB *sql.DB
}
