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

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/views/index.html")
	})

	router.Post("/sessions", api.AddSessionHandler)
	router.Get("/sessions", api.GetSessionsHandler)

	fmt.Println("Starting server at port 8080...")
	http.ListenAndServe(":8080", router)
}
