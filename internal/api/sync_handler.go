package api

import (
	"github.com/rs/zerolog/log"
	"gourd/internal/common"
	"gourd/internal/git_ops"
	"gourd/internal/middleware"
	"gourd/internal/storage"
	"net/http"
)

func (h *DBHandler) SyncProgress(w http.ResponseWriter, r *http.Request) {
	token := middleware.GetTokenFromContext(r.Context())
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	code := r.FormValue("code")
	ext := r.FormValue("mode")
	log.Info().Msgf("Submitted Code: %s", code)
	session, err := storage.GetSession(h.DB, token)
	source, err := common.GetActiveConfig().Find(session.Repo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = git_ops.CommitToBranch(session, source.LocalPath, code, ext)
	log.Info().Msg("FsMutex unlocked")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
