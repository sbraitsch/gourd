package cmd

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
	"gourd/internal/api"
	"gourd/internal/common"
	"gourd/internal/config"
	gourdMW "gourd/internal/middleware"
	"gourd/internal/storage"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Long:  `Starts the http server and begins serving server-side rendered HTML to an HTMX frontend.`,
	Run:   serve,
}

var cfgPath string

func init() {
	// add a --config param to the CLI to pass in the location of the config.toml
	serveCmd.Flags().StringVarP(&cfgPath, "config", "c", "./config.toml", "Path of the config file")
	rootCmd.AddCommand(serveCmd)
}

// serve loads the config from the config.toml and sets up the chi router, then starts the server.
func serve(cmd *cobra.Command, args []string) {
	config.LoadConfig(cfgPath)
	db := storage.ConnectDB()
	defer db.Close()
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get database driver")
	}
	var migrationsPath string
	if os.Getenv("RUNNING_TEST") == "true" {
		migrationsPath = "file://../internal/storage/migrations_test"
	} else {
		migrationsPath = "file://internal/storage/migrations"
	}
	migrations, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres", driver)
	if err != nil {
		log.Info().Err(err).Msg("Failed to get migrations instance")
	}
	if err = migrations.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}
	storage.InitAdminUser(db)

	authMW := gourdMW.AuthMiddleware{DB: db}
	dbHandler := api.HandlerStruct{DB: db}

	protectedRouter := api.ConfigureProtectedRouter(&authMW, dbHandler)
	adminRouter := api.ConfigureAdminRouter(&authMW, dbHandler)
	router := api.ConfigureMainRouter(protectedRouter, adminRouter, dbHandler)

	log.Info().Msgf("Starting server on port :%v", common.GetActiveConfig().ServerPort)
	http.ListenAndServe(fmt.Sprintf(":%v", common.GetActiveConfig().ServerPort), router)
}
