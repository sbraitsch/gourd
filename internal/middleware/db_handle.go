package middleware

import (
	"context"
	"database/sql"
	"net/http"
)

type contextKey string

const dbContextKey contextKey = "db"

func DBMiddleware(db *sql.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), dbContextKey, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetDBFromContext(ctx context.Context) *sql.DB {
	db, ok := ctx.Value(dbContextKey).(*sql.DB)
	if !ok {
		return nil
	}
	return db
}
