package git

import (
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"

	"errors"
	"log"
)

func GetBlame(repoId uint64, commitHash string, path string) ([]BlameLine, error) {
	repo, err := git.PlainOpen(realPath(repoId))
	if err != nil {
		return nil, err
	}

	objId, ok := plumbing.FromHex(commitHash)
	if !ok {
		return nil, ErrInvalidCommitHash
	}
	commit, err := repo.CommitObject(objId)
	if errors.Is(err, plumbing.ErrObjectNotFound) {
		return nil, ErrInvalidCommitHash
	} else if err != nil {
		log.Printf("Error while getting commit object: %v\n", err.Error())
		return nil, err
	}

	blameRes, err := git.Blame(commit, path)
	if errors.Is(err, object.ErrFileNotFound) {
		return nil, ErrBlobNotFound
	}

	blame := []BlameLine{}
	for _, line := range blameRes.Lines {
		blame = append(blame, BlameLine{
			Author: author{
				Name:  line.AuthorName,
				Email: line.Author,
			},
			Text:      line.Text,
			Timestamp: line.Date.Unix(),
			Hash:      line.Hash.String(),
		})
	}
	return blame, nil
}
