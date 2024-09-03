package api

import (
	"fmt"
	gourdMW "gourd/internal/middleware"
	"gourd/internal/storage"
	"gourd/internal/views"
	"log"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	token := r.FormValue("token")

	defer r.Body.Close()
	cookie := &http.Cookie{
		Name:     "Session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   3600 * 24 * 7,
	}

	w.Header().Set("HX-Trigger", "content-refresh")
	http.SetCookie(w, cookie)
	_, _ = fmt.Fprintf(w, token)
}

func AddSessionHandler(w http.ResponseWriter, r *http.Request) {
	dbConn := gourdMW.GetDBFromContext(r.Context())
	token := storage.CreateSession(dbConn, "John", "Doe", 30)
	w.Header().Set("HX-Trigger", "content-refresh")
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
