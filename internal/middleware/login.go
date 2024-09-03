package middleware

import (
	"context"
	"errors"
	"gourd/internal/views"
	"net/http"
	"strings"
)

const tokenContextKey contextKey = "token"

func LoginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")

		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				w.Header().Set("Content-Type", "text/html")
				adminRequired := strings.Contains(r.URL.Path, "/admin")
				views.LoginOverlay(adminRequired).Render(r.Context(), w)
				return
			}
			http.Error(w, "Error retrieving cookie", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), tokenContextKey, cookie.Value)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetTokenFromContext(ctx context.Context) string {
	token, ok := ctx.Value(tokenContextKey).(string)
	if !ok {
		return ""
	}
	return token
}
