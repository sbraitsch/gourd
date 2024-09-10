package common

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/url"
	"strings"
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
	URL         string
	LocalPath   string
	DisplayName string
	Username    string
	PAT         string
}

var activeConfig *Config

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

func (cfg *Config) Find(repo string) (*Source, error) {
	for _, source := range cfg.Sources {
		if source.URL == repo {
			return &source, nil
		}
	}
	return nil, fmt.Errorf("no source configured for repository %s", repo)
}

func GetActiveConfig() *Config {
	return activeConfig
}

func SetActiveConfig(cfg *Config) {
	activeConfig = cfg
}
