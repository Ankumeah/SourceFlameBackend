package database

import (
	"github.com/Ankumeah/SourceFlameBackend/internal/git"

	"context"
	"errors"
	"log"
)

func (d *GitDb) CreateRepo(
	ctx context.Context,
	ownerId uint64,
	repoName string,
) (uint64, error) {
	repoId, err := d.db.CreateRepo(ctx, ownerId, repoName)
	if errors.Is(err, ErrSafe) {
		return 0, err
	} else if err != nil {
		log.Printf("Error while adding repo: %v\n", err.Error())
		return 0, err
	}

	err = git.CreateRepo(repoId)
	if err != nil {
		log.Printf("Error while creating repo dir: %v\n", err.Error())
		d.db.DeleteRepo(ctx, repoId)

		return 0, err
	}

	return repoId, nil
}

func (d *GitDb) GetId(
	ctx context.Context,
	ownerId uint64,
	repoName string,
) (uint64, error) {
	repoId, err := d.db.GetId(ctx, ownerId, repoName)
	if !errors.Is(err, ErrInvalid) && err != nil {
		log.Printf("Error while getting repo id: %v\n", err.Error())
	}

	return repoId, err
}

func (d *GitDb) DeleteRepo(
	ctx context.Context,
	repoId uint64,
) error {
	err := d.db.DeleteRepo(ctx, repoId)
	if !errors.Is(err, ErrInvalid) && err != nil {
		log.Printf("Error while removing repo: %v\n", err.Error())
		return err
	}

	err = git.DeleteRepo(repoId)
	if err != nil {
		log.Printf("Error while deleting repo dir: %v\n", err.Error())
		return err
	}

	return nil
}

func (d *GitDb) GetRepos(
	ctx context.Context,
	userId uint64,
	limit uint8,
	offset uint64,
) ([]RepoInfo, error) {
	var maxLimit uint8 = 100
	if limit > maxLimit {
		return nil, ErrLimitTooLarge
	}

	repos, err := d.db.GetRepos(ctx, userId, limit, offset)
	if !errors.Is(err, ErrInvalid) && err != nil {
		log.Printf("Error while getting repos: %v\n", err.Error())
	}

	return repos, err
}

func (d *GitDb) Info(
	ctx context.Context,
	repoId uint64,
) (*RepoInfo, error) {
	info, err := d.db.Info(ctx, repoId)
	if !errors.Is(err, ErrInvalid) && err != nil {
		log.Printf("Error while getting repo info: %v\n", err.Error())
		return nil, err
	}

	return info, err
}

func (d *GitDb) TransferOwner(
  ctx context.Context,
  repoId uint64,
  newOwnerId uint64,
) error {
  err := d.db.TransferOwner(ctx, repoId, newOwnerId)
  if !errors.Is(err, ErrSafe) && err != nil {
		log.Printf("Error while transfering repo owner: %v\n", err.Error())
  }
  return err
}
