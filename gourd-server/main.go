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
  /*
  db := ConnectDB()
	defer db.Close()

	CreateTable(db)

	InsertPerson(db, "John Doe", "john.doe@example.com")

	people, err := GetPeople(db)
	if err != nil {
		log.Fatal("Error getting people: ", err)
	}

	for _, person := range people {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", person.ID, person.Name, person.Email)
	}
  */

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
