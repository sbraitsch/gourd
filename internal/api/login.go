package api

import (
	"encoding/json"
	"fmt"
	gourdMW "gourd/internal/middleware"
	"gourd/internal/storage"
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

	dbConn := gourdMW.GetDBFromContext(r.Context())
	session, err := storage.GetSession(dbConn, token)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Token not recognized", http.StatusNotFound)
		return
	}

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

	jsonData, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		http.Error(w, "Unable to marshal return value", http.StatusInternalServerError)
		return
	}
	_, _ = fmt.Fprintf(w, string(jsonData))
}

func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
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
