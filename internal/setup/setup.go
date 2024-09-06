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

	protectedRouter := configureProtectedRouter(&authMW, db)
	adminRouter := configureAdminRouter(&authMW, db)
	router := configureMainRouter(protectedRouter, adminRouter, db)
	return router, db
}

func configureMainRouter(protectedRouter, adminRouter *chi.Mux, db *sql.DB) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/views/index.html")
	})

	loginHandler := api.LoginHandler{DB: db}
	router.Post("/login", loginHandler.Login)
	router.Get("/clone", api.CloneHandler)

	router.Mount("/api", protectedRouter)
	router.Mount("/admin", adminRouter)
	return router
}

func configureAdminRouter(authMW *gourdMW.AuthMiddleware, db *sql.DB) *chi.Mux {
	router := chi.NewRouter()
	router.Use(authMW.Authenticate)
	router.Use(authMW.AuthenticateAdmin)

	sessionHandler := api.SessionHandler{DB: db}
	router.Post("/generate", sessionHandler.GenerateSession)
	return router
}

func configureProtectedRouter(authMW *gourdMW.AuthMiddleware, db *sql.DB) *chi.Mux {
	router := chi.NewRouter()
	router.Use(authMW.Authenticate)

	questionHandler := api.QuestionHandler{DB: db}
	router.Get("/questions", questionHandler.GetQuestion)
	contentHandler := api.ContentHandler{DB: db}
	router.Get("/content", contentHandler.GetContent)
	return router
}
