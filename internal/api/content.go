package api

import (
	"gourd/internal/middleware"
	"gourd/internal/storage"
	"gourd/internal/views"
	"net/http"
)

func ContentHandler(w http.ResponseWriter, r *http.Request) {
	isAdmin := middleware.GetAdminStatusFromContext(r.Context())
	if !isAdmin {
		session, err := storage.GetSession(middleware.GetDBFromContext(r.Context()), middleware.GetTokenFromContext(r.Context()))
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
		views.SessionGenerator().Render(r.Context(), w)
	}
}
