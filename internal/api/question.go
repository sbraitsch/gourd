package api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
	"gourd/internal/views"
	"io"
	"net/http"
	"os"
)

func QuestionHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "../gourd_example/question_01.md"

	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v", err)
	}

	md := goldmark.New()

	var buf bytes.Buffer
	if err := md.Convert(content, &buf); err != nil {
		http.Error(w, "Error rendering question source", http.StatusInternalServerError)
		return
	}

	intro := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, buf.String())
		return
	})

	views.Question(intro).Render(r.Context(), w)
}
