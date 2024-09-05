package api

import (
	"database/sql"
	"gourd/internal/config"
	"gourd/internal/storage"
	"gourd/internal/views"
	"log"
	"net/http"
	"strconv"
)

type SessionHandler struct {
	DB *sql.DB
}

func (h SessionHandler) GetSessionGenerator(w http.ResponseWriter, r *http.Request) {
	views.SessionGenerator(config.ActiveConfig.Sources).Render(r.Context(), w)
}

func (h SessionHandler) GenerateSession(w http.ResponseWriter, r *http.Request) {
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
	userId := storage.CreateUser(h.DB, firstname, lastname, false)
	token := storage.CreateSession(h.DB, userId, repo, timelimit)
	views.GenerationResult(token.String()).Render(r.Context(), w)
}

func (h SessionHandler) GetSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := storage.GetSessions(h.DB)

	if err != nil {
		log.Fatal("Error retrieving sessions from database: ", err)
	}
	views.List(sessions).Render(r.Context(), w)
}
