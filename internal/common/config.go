package common

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/url"
	"os"
	"strings"
)

// Config represents the content of the config.toml in code.
type Config struct {
	ApplicationTitle    string
	ApplicationSubtitle string
	LogoPath            string
	ServerPort          int

	DB DBConfig

	Sources []Source
}

type DBConfig struct {
	Password string
	User     string
	Name     string
	Host     string
	Port     string
}

type Source struct {
	URL         string
	LocalPath   string
	DisplayName string
	Username    string
	PAT         string
}

var activeConfig *Config

// GetRepoName returns the repositories name from its URL.
func (src *Source) GetRepoName() string {
	parsedURL, err := url.Parse(src.URL)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse url")
		return ""
	}

	parts := strings.Split(parsedURL.Path, "/")
	if len(parts) < 3 {
		return ""
	}

	return strings.TrimSuffix(parts[2], ".git")
}

// Find searches for the given repository in the list of configured sources, returning it if found.
func (cfg *Config) Find(repo string) (*Source, error) {
	for _, source := range cfg.Sources {
		if source.URL == repo {
			return &source, nil
		}
	}
	return nil, fmt.Errorf("no source configured for repository %s", repo)
}

// GetActiveConfig returns the active config.
func GetActiveConfig() *Config {
	return activeConfig
}

// SetActiveConfig sets the active config.
func SetActiveConfig(cfg *Config) {
	activeConfig = cfg
}

func IsTestEnvironment() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test") {
			return true
		}
	}
	return false
}
