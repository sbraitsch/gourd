package api

import (
	"database/sql"
	"gourd/internal/config"
	"gourd/internal/middleware"
	"gourd/internal/storage"
	"gourd/internal/views"
	"net/http"
)

type ContentHandler struct {
	DB *sql.DB
}

func (h ContentHandler) GetContent(w http.ResponseWriter, r *http.Request) {
	isAdmin := middleware.GetAdminStatusFromContext(r.Context())
	if !isAdmin {
		session, err := storage.GetSession(h.DB, middleware.GetTokenFromContext(r.Context()))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		intro, err := RenderQuestion()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		views.Question(intro, session).Render(r.Context(), w)
	} else {
		views.SessionGenerator(config.ActiveConfig.Sources).Render(r.Context(), w)
	}
}
