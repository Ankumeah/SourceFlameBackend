package git

import (
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/filemode"
	"github.com/go-git/go-git/v6/plumbing/object"

	"log"
)

func List_Dir(repo_id uint64, commit_hash string, path string) ([]File, error) {
  repo, err := git.PlainOpen(real_path(repo_id))
  if err != nil {
    return nil, err
  }

  hash := plumbing.NewHash(commit_hash)
  commit, err := repo.CommitObject(hash)
  if err == plumbing.ErrObjectNotFound {
    return nil, Error_Inavlid_Commit_Hash
  } else if err != nil {
    log.Printf("Error while getting commit object: %v\n", err.Error())
    return nil, err
  }

  tree, err := commit.Tree()
  if path != "" { tree, err = tree.Tree(path) }
  if err == object.ErrDirectoryNotFound || err == object.ErrEntryNotFound {
    return nil, Error_Path_Not_Found
  } else if err == object.ErrMaxTreeDepth {
    return nil, Error_Path_Too_Deep
  } else if err != nil {
    log.Printf("Error while getting commit worktree: %v\n", err.Error())
    return nil, err
  }

  var files []File
  entries := tree.Entries

  for _, entry := range entries {
    files = append(files, File {
      File_name: entry.Name,
      Dir: entry.Mode == filemode.Dir,
    })
  }

  return files, nil
}
