package api

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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

func (h *DBHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	token, err := parseRequestToken(w, r)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse token")
		return
	}
	exists := storage.CheckUserExists(h.DB, token, false)

	if !exists {
		log.Error().Msg("Token not recognized")
		http.Error(w, "Token not recognized", http.StatusNotFound)
		return
	}

	writeCookies(token, w)
	log.Info().Msg("Login successful")
	fmt.Fprint(w, "Login successful")
}

func parseRequestToken(w http.ResponseWriter, r *http.Request) (string, error) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return "", err
	}

	token := r.FormValue("token")
	if _, err = uuid.Parse(token); err != nil {
		fmt.Println(err)
		http.Error(w, "Malformed Input", http.StatusBadRequest)
		return "", err
	}
	return token, nil
}
