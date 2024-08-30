package api

import (
	"fmt"
	gourdMW "gourd/internal/middleware"
	"gourd/internal/storage"
	"gourd/internal/views"
	"log"
	"net/http"
)

func AddSessionHandler(w http.ResponseWriter, r *http.Request) {
	dbConn := gourdMW.GetDBFromContext(r.Context())
	token := storage.CreateSession(dbConn, "John", "Doe", 30)
	w.Header().Set("HX-Trigger", "list-refresh")
	_, _ = fmt.Fprintf(w, token.String())
}

func GetSessionsHandler(w http.ResponseWriter, r *http.Request) {
	dbConn := gourdMW.GetDBFromContext(r.Context())
	sessions, err := storage.GetSessions(dbConn)

	if err != nil {
		log.Fatal("Error retrieving sessions from database: ", err)
	}
	views.List(sessions).Render(r.Context(), w)
}
