package git

import "os"

func Delete_Repo(repo_id uint64) error {
  return os.RemoveAll(real_path(repo_id))
}
