package setup

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gourd/internal"
	"gourd/internal/api"
	"gourd/internal/common"
	gourdMW "gourd/internal/middleware"
	"gourd/internal/storage"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

func Init() (*chi.Mux, *sql.DB) {
	db := storage.ConnectDB()
	storage.CreateTable(db)

	authMW := gourdMW.AuthMiddleware{DB: db}

	dbHandler := api.DBHandler{DB: db}
	protectedRouter := configureProtectedRouter(&authMW, dbHandler)
	adminRouter := configureAdminRouter(&authMW, dbHandler)
	router := configureMainRouter(protectedRouter, adminRouter, dbHandler)
	return router, db
}

func configureMainRouter(protectedRouter, adminRouter *chi.Mux, handler api.DBHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	fs := http.FileServer(http.FS(internal.StaticAssets))
	router.Handle("/internal/static/*", http.StripPrefix("/internal/", fs))

	customLogoPath := common.GetActiveConfig().LogoPath
	relLogoDir := filepath.Dir(customLogoPath)
	routingDir := relLogoDir
	if strings.HasPrefix(routingDir, "..") {
		routingDir = routingDir[2:]
	}
	if !strings.HasPrefix(routingDir, "/") {
		routingDir = "/" + routingDir
	}

	customFs := http.FileServer(http.Dir(relLogoDir))
	router.Handle(fmt.Sprintf("%s/*", routingDir), http.StripPrefix(routingDir, customFs))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		file, err := internal.StaticAssets.Open("static/index.html")
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
		}
		defer file.Close()
		fileInfo, err := file.Stat()
		if err != nil {
			http.Error(w, "Could not get file info", http.StatusInternalServerError)
		}
		http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file.(io.ReadSeeker))
	})
	router.Post("/login", handler.Login)

	router.Mount("/api", protectedRouter)
	router.Mount("/admin", adminRouter)
	return router
}

func configureAdminRouter(authMW *gourdMW.AuthMiddleware, handler api.DBHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Use(authMW.Authenticate)
	router.Use(authMW.AuthenticateAdmin)

	router.Post("/generate", handler.GenerateSession)
	return router
}

func configureProtectedRouter(authMW *gourdMW.AuthMiddleware, handler api.DBHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Use(authMW.Authenticate)

	router.Get("/questions/{id}", handler.GetQuestion)
	router.Get("/content", handler.GetContent)
	router.Post("/sync", handler.SyncProgress)
	return router
}
