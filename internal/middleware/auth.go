package middleware

import (
	"context"
	"database/sql"
	"errors"
	"github.com/rs/zerolog/log"
	"gourd/internal/common"
	"gourd/internal/storage"
	"gourd/internal/views"
	"net/http"
)

type contextKey string

const tokenContextKey contextKey = "token"
const adminContextKey contextKey = "admin"

type AuthMiddleware struct {
	DB *sql.DB
}

// AuthenticationBasic is router middleware that authenticates requests based on a token cookie.
// If the cookie does not exist, it returns a rendered Login-HTML.
func (mw *AuthMiddleware) AuthenticationBasic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")

		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				w.Header().Set("Content-Type", "text/html")
				cfg := common.GetActiveConfig()
				views.Login(cfg.ApplicationTitle, cfg.ApplicationSubtitle, cfg.LogoPath).Render(r.Context(), w)
				return
			}
			http.Error(w, "Error retrieving cookie", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), tokenContextKey, cookie.Value)
		exists, isAdmin := storage.CheckUserExists(mw.DB, cookie.Value)
		if !exists {
			log.Error().Msg("Token not recognized")
			http.Error(w, "Token not recognized", http.StatusNotFound)
			return
		}
		ctx = context.WithValue(ctx, adminContextKey, isAdmin)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthenticationAdmin is an additional router middleware for admin endpoints.
// This is run after the basic authentication middleware, so the context is set.
// Checks if the token belongs to an admin user.
func (mw *AuthMiddleware) AuthenticationAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin := GetAdminStatusFromContext(r.Context())
		if !isAdmin {
			http.Error(w, "Missing privileges", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetTokenFromContext retrieves the token from the request context.
func GetTokenFromContext(ctx context.Context) string {
	token, ok := ctx.Value(tokenContextKey).(string)
	if !ok {
		return ""
	}
	return token
}

// GetAdminStatusFromContext retrieves the admin status from the request context.
func GetAdminStatusFromContext(ctx context.Context) bool {
	isAdmin, ok := ctx.Value(adminContextKey).(bool)
	if !ok {
		return false
	}
	return isAdmin
}
