package git

import "github.com/go-git/go-git/v6"

func Create_Repo(repo_id uint64) error {
	_, err := git.PlainInit(real_path(repo_id), true)
	if err == git.ErrTargetDirNotEmpty {
		return Error_Repository_Exists
	} else if err != nil {
		return err
	}

	return nil
}
