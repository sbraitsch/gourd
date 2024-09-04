package api

import (
	gourdMW "gourd/internal/middleware"
	"gourd/internal/storage"
	"gourd/internal/views"
	"log"
	"net/http"
	"strconv"
)

func SessionGeneratorHandler(w http.ResponseWriter, r *http.Request) {
	views.SessionGenerator().Render(r.Context(), w)
}

func GenerateSessionHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	firstname := r.FormValue("firstname")
	lastname := r.FormValue("lastname")
	timelimit, err := strconv.ParseInt(r.FormValue("timelimit"), 10, 64)
	if err != nil {
		log.Fatal("Someone managed to fuck up the time limit input.")
	}
	repo := r.FormValue("repo")

	dbConn := gourdMW.GetDBFromContext(r.Context())
	userId := storage.CreateUser(dbConn, firstname, lastname, false)
	token := storage.CreateSession(dbConn, userId, repo, timelimit)
	views.GenerationResult(token.String()).Render(r.Context(), w)
}

func GetSessionsHandler(w http.ResponseWriter, r *http.Request) {
	dbConn := gourdMW.GetDBFromContext(r.Context())
	sessions, err := storage.GetSessions(dbConn)

	if err != nil {
		log.Fatal("Error retrieving sessions from database: ", err)
	}
	views.List(sessions).Render(r.Context(), w)
}
