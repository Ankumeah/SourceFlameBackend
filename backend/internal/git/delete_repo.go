package git

import "os"

func DeleteRepo(repoId uint64) error {
	return os.RemoveAll(realPath(repoId))
}
