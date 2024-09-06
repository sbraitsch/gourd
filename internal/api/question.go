package api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
	"gourd/internal/middleware"
	"gourd/internal/storage"
	"gourd/internal/views"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func (h DBHandler) GetQuestion(w http.ResponseWriter, r *http.Request) {
	token := middleware.GetTokenFromContext(r.Context())
	session, err := storage.GetSession(h.DB, token)
	intro, code, mode, err := RenderQuestion(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	views.Question(intro, code, mode, session).Render(r.Context(), w)
}

func RenderQuestion(session storage.Session) (templ.Component, string, string, error) {
	questionFilePath := fmt.Sprintf("../gourd_example/part_%02d/question.md", session.Step)
	providedCodeFilePath := fmt.Sprintf("../gourd_example/part_%02d/provided_code.*", session.Step)

	question, err := os.ReadFile(questionFilePath)
	if err != nil {
		return nil, "", "", err
	}

	md := goldmark.New()

	var buf bytes.Buffer
	if err := md.Convert(question, &buf); err != nil {
		return nil, "", "", err
	}

	paths, err := filepath.Glob(providedCodeFilePath)
	if err != nil {
		return nil, "", "", err
	}
	code, err := os.ReadFile(paths[0])
	if err != nil {
		return nil, "", "", err
	}
	ext := filepath.Ext(paths[0])

	intro := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, buf.String())
		return
	})

	return intro, string(code), resolveMode(ext), nil
}

func resolveMode(ext string) string {
	switch ext {
	case ".go":
		return "go"
	case ".java":
		return "text/x-java"
	case ".rs":
		return "rust"
	case ".ts":
		return "application/typescript"
	case ".js":
		return "javascript"
	case ".py":
		return "python"
	case ".kt":
		return "text/x-kotlin"
	default:
		return "text/plain" // Default mode if extension is unknown
	}
}
