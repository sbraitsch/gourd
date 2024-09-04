package api

import (
	"bytes"
	"context"
	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
	"gourd/internal/middleware"
	"gourd/internal/storage"
	"gourd/internal/views"
	"io"
	"net/http"
	"os"
)

func QuestionHandler(w http.ResponseWriter, r *http.Request) {
	intro, err := RenderQuestion()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session, err := storage.GetSession(middleware.GetDBFromContext(r.Context()), middleware.GetTokenFromContext(r.Context()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	views.Question(intro, session).Render(r.Context(), w)
}

func RenderQuestion() (templ.Component, error) {
	filePath := "../gourd_example/question_01.md"

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	md := goldmark.New()

	var buf bytes.Buffer
	if err := md.Convert(content, &buf); err != nil {
		return nil, err
	}

	intro := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, buf.String())
		return
	})

	return intro, nil
}
