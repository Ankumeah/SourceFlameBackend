package git

import (
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"

	"log"
)

func GetBranches(repoId uint64) ([]string, error) {
	repo, err := git.PlainOpen(realPath(repoId))
	if err != nil {
		return nil, err
	}

	branchesIter, err := repo.Branches()
	if err != nil {
		log.Printf("Error while getting repo branches: %v\n", err.Error())
		return nil, err
	}
	defer branchesIter.Close()

	branches := []string{}
	branchesIter.ForEach(func(r *plumbing.Reference) error {
		branches = append(branches, r.Name().Short())
		return nil
	})

	return branches, nil
}
