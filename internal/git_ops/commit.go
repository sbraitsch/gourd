package git_ops

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/rs/zerolog/log"
	"gourd/internal/common"
	"gourd/internal/storage"
	"os"
	"sync"
	"time"
)

var fsMutex sync.Mutex

func CommitToBranch(session storage.HydratedSession, repoPath string, code, ext string) error {
	log.Info().Msg("Acquiring FS Mutex")
	fsMutex.Lock()
	defer fsMutex.Unlock()
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Info().Msgf("Couldn't open repo at %s: %v", repoPath, err)
		return err
	}
	log.Info().Msg("FS Mutex locked")
	err = CheckoutBranch(repo, session.User.GetBranchName())
	if err != nil {
		log.Info().Msgf("Couldn't checkout user branch: %v", err)
		return err
	}
	defer CheckoutBranch(repo, "main")
	err = commit(repo, repoPath, 1, code, ext)
	if err != nil {
		return err
	}
	log.Info().Msg("FS Mutex unlocked")
	return nil
}

func commit(repo *git.Repository, repoPath string, part int, content, mode string) error {
	fileName := fmt.Sprintf("part_%02d/code%s", part, common.ResolveExtMode(mode))
	filePath := fmt.Sprintf("%s/%s", repoPath, fileName)
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}
	_, err = worktree.Add(fileName)
	if err != nil {
		return err
	}
	commitMessage := fmt.Sprintf("Committing file %s", fileName)
	_, err = worktree.Commit(commitMessage, &git.CommitOptions{
		// adjust later
		Author: &object.Signature{
			Name:  "Your Name",
			Email: "your-email@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
