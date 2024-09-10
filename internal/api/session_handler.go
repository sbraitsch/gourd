package api

import (
	"github.com/go-git/go-git/v5"
	"github.com/rs/zerolog/log"
	"gourd/internal/common"
	"gourd/internal/git_ops"
	"gourd/internal/storage"
	"gourd/internal/views"
	"net/http"
	"strconv"
)

func (h *DBHandler) GenerateSession(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	firstname := r.FormValue("firstname")
	lastname := r.FormValue("lastname")
	timelimit, err := strconv.ParseInt(r.FormValue("timelimit"), 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing timelimit")
	}
	repo := r.FormValue("repo")
	user := storage.CreateUser(h.DB, firstname, lastname, false)
	storage.CreateSession(h.DB, user.ID, repo, timelimit)
	// create git branch
	source, err := common.GetActiveConfig().Find(repo)
	if err != nil {
		log.Error().Err(err).Msg("Error finding local repo configuration")
	}
	localRepo, err := git.PlainOpen(source.LocalPath)
	if err != nil {
		log.Error().Err(err).Msg("Unable to open local repository")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = git_ops.CreateBranch(localRepo, user)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create branch")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	views.GenerationResult(user.ID.String()).Render(r.Context(), w)
}
