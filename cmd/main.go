package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gourd/internal/api"
	gourdMW "gourd/internal/middleware"
	"gourd/internal/storage"
	"net/http"
)

func main() {
	db := storage.ConnectDB()
	defer db.Close()
	storage.CreateTable(db)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(gourdMW.DBMiddleware(db))

	protectedRouter := chi.NewRouter()
	protectedRouter.Use(gourdMW.AuthMiddleware)

	adminRouter := chi.NewRouter()
	adminRouter.Use(gourdMW.AuthMiddleware)
	adminRouter.Use(gourdMW.AdminMiddleware)

	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/views/index.html")
	})

	router.Post("/login", api.LoginHandler)
	router.Get("/clone", api.CloneHandler)

	protectedRouter.Get("/questions", api.QuestionHandler)
	protectedRouter.Get("/content", api.ContentHandler)

	adminRouter.Get("/sessions", api.GetSessionsHandler)
	adminRouter.Get("/generator", api.SessionGeneratorHandler)
	adminRouter.Post("/generate", api.GenerateSessionHandler)

	router.Mount("/api", protectedRouter)
	router.Mount("/admin", adminRouter)

	fmt.Println("Starting server at port 8080...")
	http.ListenAndServe(":8080", router)
}
