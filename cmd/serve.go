package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gourd/internal/config"
	"gourd/internal/setup"
	"net/http"

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
	serveCmd.Flags().StringVarP(&cfgPath, "config", "c", "./", "Path to the config file")
	rootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) {
	config.LoadConfig(cfgPath)
	router, db := setup.Init()
	defer db.Close()
	log.Info().Msgf("Starting server on port :%v", config.ActiveConfig.ServerPort)
	http.ListenAndServe(fmt.Sprintf(":%v", config.ActiveConfig.ServerPort), router)
}
