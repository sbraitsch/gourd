package git_ops

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/rs/zerolog/log"
	"gourd/internal/storage"
)

func CreateBranch(repo *git.Repository, user storage.User) error {
	headRef, err := repo.Head()
	if err != nil {
		return err
	}

	branchName := user.GetBranchName()
	newBranchRef := plumbing.NewHashReference(plumbing.NewBranchReferenceName(branchName), headRef.Hash())
	err = repo.Storer.SetReference(newBranchRef)

	if err != nil {
		return err
	}

	log.Info().Msgf("Created branch %s", branchName)
	return nil
}
