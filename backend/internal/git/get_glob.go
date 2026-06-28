package git

import (
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"

	"log"
)

func Get_Glob(repo_id uint64, commit_hash string, path string) (string, error) {
	repo, err := git.PlainOpen(real_path(repo_id))
	if err != nil {
		return "", err
	}

	hash := plumbing.NewHash(commit_hash)
	commit, err := repo.CommitObject(hash)
	if err == plumbing.ErrObjectNotFound {
		return "", Error_Inavlid_Commit_Hash
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
	if err == object.ErrDirectoryNotFound {
		return "", Error_Path_Not_Found
	} else if err == object.ErrFileNotFound {
		return "", Error_Blob_Not_Found
	} else if err != nil {
		log.Printf("Error while opening blob file: %v\n", err.Error())
		return "", err
	}

	if file.Size >= max_blob_size {
		return "", Error_Blob_Too_Large
	}

	return file.Contents()
}
