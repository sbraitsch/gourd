package git_ops

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/rs/zerolog/log"
	"gourd/internal/common"
	"os"
)

func TryClone(source common.Source) {
	log.Info().Msgf("Preparing to clone %s to %s", source.DisplayName, source.LocalPath)
	if !directoryExists(source.LocalPath) {
		options := git.CloneOptions{URL: source.URL, Progress: os.Stdout, Auth: &http.BasicAuth{Username: source.Username, Password: source.PAT}}
		_, err := git.PlainClone(source.LocalPath, false, &options)
		if err != nil {
			log.Error().Msgf("%e", err)
		}
	} else {
		log.Info().Msgf("Local repository for %s already exists, skipping clone", source.DisplayName)
	}
}

func directoryExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
