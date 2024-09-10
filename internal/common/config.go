package common

import (
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
