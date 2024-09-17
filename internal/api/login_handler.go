package api

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gourd/internal/storage"
	"net/http"
)

/*
	setCookie configures the cookie to be set on the client.

Also adds the HX-Trigger header to trigger a content refresh since the user will now be authenticated.
*/
func setCookie(token string, w http.ResponseWriter) {
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

// Login is the HandlerFunc for the /login endpoint. Parses the token from the form and checks if that user exists.
func (h *HandlerStruct) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	token, err := parseRequestForm(w, r)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse token")
		return
	}
	exists, _ := storage.CheckUserExists(h.DB, token)

	if !exists {
		log.Error().Msg("Token not recognized")
		http.Error(w, "Token not recognized", http.StatusNotFound)
		return
	}

	setCookie(token, w)
	log.Info().Msg("Login successful")
	fmt.Fprint(w, "Login successful")
}

// parseRequestForm parses the form and asserts the token is a valid UUID.
func parseRequestForm(w http.ResponseWriter, r *http.Request) (string, error) {
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
