package setup

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gourd/internal/api"
	gourdMW "gourd/internal/middleware"
	"gourd/internal/storage"
	"net/http"
)

func Init() (*chi.Mux, *sql.DB) {
	db := storage.ConnectDB()
	storage.CreateTable(db)

	authMW := gourdMW.AuthMiddleware{DB: db}

	dbHandler := api.DBHandler{DB: db}
	protectedRouter := configureProtectedRouter(&authMW, dbHandler)
	adminRouter := configureAdminRouter(&authMW, dbHandler)
	router := configureMainRouter(protectedRouter, adminRouter, dbHandler)
	return router, db
}

func configureMainRouter(protectedRouter, adminRouter *chi.Mux, handler api.DBHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/views/index.html")
	})
	router.Post("/login", handler.Login)

	router.Mount("/api", protectedRouter)
	router.Mount("/admin", adminRouter)
	return router
}

func configureAdminRouter(authMW *gourdMW.AuthMiddleware, handler api.DBHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Use(authMW.Authenticate)
	router.Use(authMW.AuthenticateAdmin)

	router.Post("/generate", handler.GenerateSession)
	return router
}

func configureProtectedRouter(authMW *gourdMW.AuthMiddleware, handler api.DBHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Use(authMW.Authenticate)

	router.Get("/questions", handler.GetQuestion)
	router.Get("/content", handler.GetContent)
	return router
}
