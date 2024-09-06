package api

import (
	"fmt"
	"gourd/internal/config"
	"gourd/internal/middleware"
	"gourd/internal/storage"
	"gourd/internal/views"
	"net/http"
)

func (h DBHandler) GetContent(w http.ResponseWriter, r *http.Request) {
	isAdmin := middleware.GetAdminStatusFromContext(r.Context())
	if !isAdmin {
		token := middleware.GetTokenFromContext(r.Context())
		session, err := storage.GetSession(h.DB, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		intro, code, mode, err := RenderQuestion(session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		views.Question(intro, code, mode, session).Render(r.Context(), w)
	} else {
		sessions, err := storage.GetSessions(h.DB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, s := range sessions {
			fmt.Println(s)
		}
		views.SessionGenerator(views.SessionList(sessions), config.ActiveConfig.Sources).Render(r.Context(), w)
	}
}
