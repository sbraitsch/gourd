package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"sync"
)

var (
	ActiveConfig *Config
	mutex        sync.RWMutex
)

type Config struct {
	DbPassword string
	DbUser     string
	DbName     string
	DbHost     string
	DbPort     string

	ServerPort int

	Sources []Source
}

type Source struct {
	URL  string
	Name string
	PAT  string
}

func GetConfig() *Config {
	mutex.RLock()
	defer mutex.RUnlock()
	return ActiveConfig
}

func setConfig(cfg *Config) {
	mutex.Lock()
	defer mutex.Unlock()
	ActiveConfig = cfg
}

func LoadConfig(cfgPath string) {
	var cfg Config
	loadLocalConfig(cfgPath)
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info().Msgf("Config file changed: %s", e.Name)
	})
	viper.WatchConfig()

	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to decode config into struct")
	}
	setConfig(&cfg)
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
