package vcs

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type GitVCS struct {
}

func (s *GitVCS) Clone(options Options) error {
	if _, err := fmt.Fprintf(
		options.Progress,
		"cloning repository %s at branch %s to directory %s\n",
		options.RepositoryURL,
		options.Branch,
		options.Directory,
	); err != nil {
		return errors.Wrap(err, "error writing initial progress to progress writer")
	}

	repo, err := git.PlainClone(options.Directory, false, &git.CloneOptions{
		URL:        options.RepositoryURL,
		RemoteName: options.Branch,
		Progress:   options.Progress,
	})
	if err != nil {
		return errors.Wrap(
			err,
			fmt.Sprintf("error cloning repository %s at branch %s to directory %s", options.RepositoryURL, options.Branch, options.Directory),
		)
	}

	if options.Revision == "" {
		return nil
	}

	workTree, err := repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "error parsing repository work tree")
	}

	return workTree.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(options.Revision),
	})
}
