package api

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	gourdMW "gourd/internal/middleware"
	"gourd/internal/storage"
	"net/http"
)

func writeCookies(token string, w http.ResponseWriter) {
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
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	err, token, dbConn, done := parseRequest(w, r)
	if done {
		return
	}
	exists := storage.CheckUserExists(dbConn, token, false)

	if !exists {
		fmt.Println(err)
		http.Error(w, "Token not recognized", http.StatusNotFound)
		return
	}

	writeCookies(token, w)
	fmt.Fprint(w, "Login successful")
}

func parseRequest(w http.ResponseWriter, r *http.Request) (error, string, *sql.DB, bool) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return nil, "", nil, true
	}

	token := r.FormValue("token")
	if _, err = uuid.Parse(token); err != nil {
		fmt.Println(err)
		http.Error(w, "Malformed Input", http.StatusBadRequest)
		return nil, "", nil, true
	}

	dbConn := gourdMW.GetDBFromContext(r.Context())
	return err, token, dbConn, false
}
