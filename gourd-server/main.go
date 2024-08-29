package main

import (
  "fmt"
  "net/http"

  "github.com/go-chi/chi/v5"
  "github.com/go-chi/chi/v5/middleware"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "<p>Hello, Gourd speaks HTMX!</p>")
}

func main() {
  r := chi.NewRouter()
  r.Use(middleware.Logger)
  r.Use(middleware.Recoverer)

  r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "../gourd-web/index.html")
  })

  r.Get("/hello", helloHandler)

  fmt.Println("Starting server at port 8080...")
  http.ListenAndServe(":8080", r)
}
