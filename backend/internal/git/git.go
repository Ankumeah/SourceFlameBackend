package git

import (
  "github.com/go-git/go-git/v6"

  "os"
)

const base_path = "/srv/git/"

func Create_Repo(username string, repo_name string, private bool) error {
  _, err := git.PlainInit(base_path + username + repo_name, true)
  return err
}

func Delete_Repo(username string, repo_name string) error {
  return os.RemoveAll(base_path + username + repo_name)
}
