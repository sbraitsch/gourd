package api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/yuin/goldmark"
	"gourd/internal/config"
	"gourd/internal/middleware"
	"gourd/internal/storage"
	"gourd/internal/views"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func (h DBHandler) GetQuestion(w http.ResponseWriter, r *http.Request) {
	token := middleware.GetTokenFromContext(r.Context())
	session, err := storage.GetSession(h.DB, token)
	idParam, err := strconv.Atoi(chi.URLParam(r, "id"))
	session.CurrentStep = idParam
	if idParam > session.MaxProgress {
		session.MaxProgress = idParam
	}
	err = storage.UpdateSessionProgress(h.DB, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	intro, code, mode, err := RenderQuestion(idParam, session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	views.Question(intro, code, mode, idParam).Render(r.Context(), w)
}

func RenderQuestion(step int, session storage.Session) (templ.Component, string, string, error) {
	var repoPath string
	for _, source := range config.ActiveConfig.Sources {
		if source.URL == session.Repo {
			repoPath = source.LocalPath
			break
		}
	}

	questionFilePath := fmt.Sprintf("%s/part_%02d/question.md", repoPath, step)
	providedCodeFilePath := fmt.Sprintf("%s/part_%02d/provided_code.*", repoPath, step)

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
