package main

import (
  "fmt",
  "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello, Gourd Server!")
}

func main() {
  http.HandleFunc("/", handler)
  fmt.Println("Starting server at port 8080...")
  http.ListenAndServer(":8080", nil)
}
