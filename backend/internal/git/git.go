package git

import (
  "github.com/go-git/go-git/v6"

  "os"
  "fmt"
)

const base_path = "/srv/git/"

func Create_Repo(repo_id uint64, private bool) error {
  path := base_path + fmt.Sprintf("%v", repo_id)
  _, err := git.PlainInit(path, true)
  if err == git.ErrTargetDirNotEmpty {
    return Error_repository_exists
  }

  return err
}

func Delete_Repo(repo_id uint64) error {
  path := base_path + fmt.Sprintf("%v", repo_id)
  return os.RemoveAll(path)
}
