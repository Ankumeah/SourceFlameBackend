package git

import (
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"

	"errors"
	"log"
)

func Get_Commits(repo_id uint64, branch string) ([]Commit, error) {
	repo, err := git.PlainOpen(real_path(repo_id))
	if err != nil {
		return nil, err
	}

	hash, err := repo.ResolveRevision(plumbing.Revision(branch))
	if errors.Is(err, plumbing.ErrReferenceNotFound) {
		return nil, Error_Branch_Not_Found
	} else if err != nil {
		log.Printf("Error while resolveing repo branch: %v\n", err.Error())
		return nil, err
	}

	commitsIter, err := repo.Log(&git.LogOptions{From: *hash})
	if err != nil {
		log.Printf("Error while getting branch commits: %v\n", err.Error())
		return nil, err
	}
	defer commitsIter.Close()

	commits := []Commit{}
	commitsIter.ForEach(func(commit *object.Commit) error {
		commits = append(commits, Commit{
			Author: author{
				Name:  commit.Author.Name,
				Email: commit.Author.Email,
			},
			Message:   commit.Message,
			Hash:      commit.Hash.String(),
			Timestamp: uint64(commit.Author.When.Unix()),
		})
		return nil
	})

	return commits, nil
}
