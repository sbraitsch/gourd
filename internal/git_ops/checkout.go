package git_ops

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/rs/zerolog/log"
	"gourd/internal/storage"
)

func CheckoutBranch(repo *git.Repository, user storage.User) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	branchName := fmt.Sprintf("%s_%s_%s", user.Firstname, user.Lastname, user.ID)
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

func CheckoutMain(repo *git.Repository) error {
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	mainBranchRef := plumbing.NewBranchReferenceName("main")

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: mainBranchRef,
		Create: false,
	})
	if err != nil {
		return err
	}

	log.Info().Msg("Checked out main")
	return nil
}
