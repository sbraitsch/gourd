package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"gourd/internal/common"
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
	// add a --config param to the CLI to pass in the location of the config.toml
	serveCmd.Flags().StringVarP(&cfgPath, "config", "c", "./", "Path to the config file")
	rootCmd.AddCommand(serveCmd)
}

// serve loads the config from the config.toml and sets up the chi router, then starts the server.
func serve(cmd *cobra.Command, args []string) {
	config.LoadConfig(cfgPath)
	router, db := setup.Init()
	defer db.Close()
	log.Info().Msgf("Starting server on port :%v", common.GetActiveConfig().ServerPort)
	http.ListenAndServe(fmt.Sprintf(":%v", common.GetActiveConfig().ServerPort), router)
}
