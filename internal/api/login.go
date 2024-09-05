package api

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"gourd/internal/storage"
	"net/http"
)

type LoginHandler struct {
	DB *sql.DB
}

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

func (h LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	err, token, done := parseRequest(w, r)
	if done {
		return
	}
	exists := storage.CheckUserExists(h.DB, token, false)

	if !exists {
		fmt.Println(err)
		http.Error(w, "Token not recognized", http.StatusNotFound)
		return
	}

	writeCookies(token, w)
	fmt.Fprint(w, "Login successful")
}

func parseRequest(w http.ResponseWriter, r *http.Request) (error, string, bool) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return nil, "", true
	}

	token := r.FormValue("token")
	if _, err = uuid.Parse(token); err != nil {
		fmt.Println(err)
		http.Error(w, "Malformed Input", http.StatusBadRequest)
		return nil, "", true
	}
	return err, token, false
}
