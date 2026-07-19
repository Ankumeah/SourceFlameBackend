package git

import "github.com/go-git/go-git/v6"

func CreateRepo(repoId uint64) error {
	_, err := git.PlainInit(realPath(repoId), true)
	if err == git.ErrTargetDirNotEmpty {
		return ErrRepositoryExists
	} else if err != nil {
		return err
	}

	return nil
}
