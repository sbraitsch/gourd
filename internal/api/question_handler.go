package api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/yuin/goldmark"
	"gourd/internal/common"
	"gourd/internal/middleware"
	"gourd/internal/storage"
	"gourd/internal/views"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// GetQuestion is the HandlerFunc for the /api/question endpoint.
func (h *DBHandler) GetQuestion(w http.ResponseWriter, r *http.Request) {
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

// RenderQuestion renders the dynamic part of the question HTML based on the provided context data.
func RenderQuestion(step int, session storage.HydratedSession) (templ.Component, string, string, error) {
	source, err := common.GetActiveConfig().Find(session.Repo)
	if err != nil {
		fmt.Errorf("repository %s not configured locally", session.Repo)
	}

	// build local file paths
	questionFilePath := fmt.Sprintf("%s/part_%02d/question.md", source.LocalPath, step)
	providedCodeFilePath := fmt.Sprintf("%s/part_%02d/code.*", source.LocalPath, step)

	question, err := os.ReadFile(questionFilePath)
	if err != nil {
		return nil, "", "", err
	}

	md := goldmark.New()

	// read the file content into markdown
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

	return intro, string(code), common.ResolveExtMode(ext), nil
}
