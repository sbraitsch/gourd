package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gourd/internal/common"
	"gourd/internal/git_ops"
)

func LoadConfig(cfgPath string) {
	var cfg common.Config
	loadLocalConfig(cfgPath)
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info().Msgf("Config file changed: %s", e.Name)
		common.GetActiveConfig().Sources = nil
		err := viper.Unmarshal(&cfg)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to decode config into struct")
		}
		common.SetActiveConfig(&cfg)
		log.Info().Msgf("Loaded config: %+v", cfg)
	})
	viper.WatchConfig()

	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to decode config into struct")
	}
	for _, source := range cfg.Sources {
		git_ops.TryClone(source)
	}
	common.SetActiveConfig(&cfg)
}

func loadLocalConfig(cfgPath string) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(cfgPath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Error().Err(err).Msgf("Error loading config")
	}

}
