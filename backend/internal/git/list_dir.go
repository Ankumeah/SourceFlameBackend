package git

import (
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/filemode"
	"github.com/go-git/go-git/v6/plumbing/object"

	"log"
  "errors"
)

func ListDir(repoId uint64, commitHash string, path string) ([]File, error) {
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

	tree, err := commit.Tree()
	if path != "" {
		tree, err = tree.Tree(path)
	}
	if errors.Is(err, object.ErrDirectoryNotFound) || errors.Is(err, object.ErrEntryNotFound) {
		return nil, ErrPathNotFound
	} else if errors.Is(err, object.ErrMaxTreeDepth) {
		return nil, ErrPathTooDeep
	} else if err != nil {
		log.Printf("Error while getting commit worktree: %v\n", err.Error())
		return nil, err
	}

	var files []File
	entries := tree.Entries

	for _, entry := range entries {
		files = append(files, File{
			FileName: entry.Name,
			Dir:      entry.Mode == filemode.Dir,
		})
	}

	return files, nil
}
