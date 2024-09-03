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
	protectedRouter.Use(gourdMW.LoginMiddleware)

	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/views/index.html")
	})

	router.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/views/admin.html")
	})

	router.Post("/login", api.LoginHandler)

	protectedRouter.Get("/sessions", api.GetSessionsHandler)
	protectedRouter.Post("/sessions", api.AddSessionHandler)

	router.Mount("/api", protectedRouter)

	fmt.Println("Starting server at port 8080...")
	http.ListenAndServe(":8080", router)
}
