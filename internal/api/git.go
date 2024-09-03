package api

import (
	"fmt"
	gitops "gourd/internal/git"
	"net/http"
)

func CloneHandler(w http.ResponseWriter, r *http.Request) {
	result := gitops.Clone("https://github.com/sbraitsch/gourd_example.git", "READMEFROMENV")
	_, _ = fmt.Fprintf(w, result)
}
