package middleware

import (
	"context"
	"errors"
	"gourd/internal/storage"
	"gourd/internal/views"
	"net/http"
)

const tokenContextKey contextKey = "token"
const adminContextKey contextKey = "admin"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")

		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				w.Header().Set("Content-Type", "text/html")
				views.LoginOverlay().Render(r.Context(), w)
				return
			}
			http.Error(w, "Error retrieving cookie", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), tokenContextKey, cookie.Value)
		isAdmin := storage.CheckUserExists(GetDBFromContext(r.Context()), cookie.Value, true)
		ctx = context.WithValue(ctx, adminContextKey, isAdmin)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin := GetAdminStatusFromContext(r.Context())
		if !isAdmin {
			http.Error(w, "Missing privileges", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetTokenFromContext(ctx context.Context) string {
	token, ok := ctx.Value(tokenContextKey).(string)
	if !ok {
		return ""
	}
	return token
}

func GetAdminStatusFromContext(ctx context.Context) bool {
	isAdmin, ok := ctx.Value(adminContextKey).(bool)
	if !ok {
		return false
	}
	return isAdmin
}
