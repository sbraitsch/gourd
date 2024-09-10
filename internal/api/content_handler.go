package api

import (
	"gourd/internal/common"
	"gourd/internal/middleware"
	"gourd/internal/storage"
	"gourd/internal/views"
	"net/http"
)

func (h *DBHandler) GetContent(w http.ResponseWriter, r *http.Request) {
	isAdmin := middleware.GetAdminStatusFromContext(r.Context())
	if !isAdmin {
		token := middleware.GetTokenFromContext(r.Context())
		session, err := storage.GetSession(h.DB, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		intro, code, mode, err := RenderQuestion(session.MaxProgress, session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		renderedQuestion := views.Question(intro, code, mode, session.MaxProgress)
		views.QuestionContainer(renderedQuestion).Render(r.Context(), w)
	} else {
		sessions, err := storage.GetSessions(h.DB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		views.SessionGenerator(views.SessionList(sessions), common.GetActiveConfig().Sources).Render(r.Context(), w)
	}
}
