package git

import (
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"

	"log"
  "errors"
)

func GetBlob(repoId uint64, commitHash string, path string) (string, error) {
	repo, err := git.PlainOpen(realPath(repoId))
	if err != nil {
		return "", err
	}

	hash := plumbing.NewHash(commitHash)
	commit, err := repo.CommitObject(hash)
	if errors.Is(err, plumbing.ErrObjectNotFound) {
		return "", ErrInvalidCommitHash
	} else if err != nil {
		log.Printf("Error while getting commit object: %v\n", err.Error())
		return "", err
	}

	tree, err := commit.Tree()
	if err != nil {
		log.Printf("Error while getting commit worktree: %v\n", err.Error())
		return "", err
	}

	file, err := tree.File(path)
	if errors.Is(err, object.ErrDirectoryNotFound) {
		return "", ErrPathNotFound
	} else if errors.Is(err, object.ErrFileNotFound) {
		return "", ErrBlobNotFound
	} else if err != nil {
		log.Printf("Error while opening blob file: %v\n", err.Error())
		return "", err
	}

	if file.Size >= maxBlobSize {
		return "", ErrBlobTooLarge
	}

	return file.Contents()
}
