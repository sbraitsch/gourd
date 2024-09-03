package api

import (
	"fmt"
	gourdMW "gourd/internal/middleware"
	"gourd/internal/storage"
	"net/http"
)

func writeCookies(token string, isAdmin bool, w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   3600 * 24 * 7,
	}

	w.Header().Set("HX-Trigger", "content-refresh")
	http.SetCookie(w, cookie)

	if isAdmin {
		adminCookie := &http.Cookie{
			Name:     "isAdmin",
			Value:    "true",
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   3600 * 24 * 7,
		}
		http.SetCookie(w, adminCookie)
	}
}

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

	writeCookies(session.ID.String(), false, w)
}

func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	token := r.FormValue("token")

	dbConn := gourdMW.GetDBFromContext(r.Context())
	isAdmin, err := storage.CheckAdminStatus(dbConn, token)

	if !isAdmin {
		fmt.Println(err)
		http.Error(w, "Token not recognized", http.StatusNotFound)
		return
	}
	defer r.Body.Close()
	writeCookies(token, true, w)
}
