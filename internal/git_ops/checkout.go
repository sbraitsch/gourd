package git_ops

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/rs/zerolog/log"
)

// CheckoutBranch swaps to the given branch of the local repository.
func CheckoutBranch(repo *git.Repository, branchName string) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	newBranchRef := plumbing.NewBranchReferenceName(branchName)

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: newBranchRef,
		Create: false,
	})
	if err != nil {
		log.Error().Err(err).Msgf("error checking out branch %s", branchName)
		return err
	}

	log.Info().Msgf("Checked out branch %s", branchName)
	return nil
}
