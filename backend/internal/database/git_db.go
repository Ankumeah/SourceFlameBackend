package database

import (
	"github.com/Ankumeah/SourceFlameBackend/internal/git"

	"context"
	"errors"
	"log"
)

func (d *Git_db) Create_Repo(
  ctx context.Context,
  owner_id uint64,
  repo_name string,
  private bool,
) (uint64, error) {
  repo_id, err := d.db.Create_Repo(ctx, owner_id, repo_name, private)
  if errors.Is(err, Safe_Error) {
    return 0, err
  } else if err != nil {
    log.Printf("Error while adding repo: %v\n", err.Error())
    return 0, err
  }

  err = git.Create_Repo(repo_id)
  if err != nil {
    log.Printf("Error while creating repo dir: %v\n", err.Error())
    d.db.Delete_Repo(ctx, repo_id)

    return 0, err
  }

  return repo_id, nil
}

func (d *Git_db) Get_Id(
  ctx context.Context,
  owner_id uint64,
  repo_name string,
) (uint64, error) {
  repo_id, err := d.db.Get_Id(ctx, owner_id, repo_name)
  if !errors.Is(err, Error_Invalid) && err != nil {
    log.Printf("Error while getting repo id: %v\n", err.Error())
  }

  return repo_id, err
}

func (d *Git_db) Delete_Repo(
  ctx context.Context,
  repo_id uint64,
) error {
  err := d.db.Delete_Repo(ctx, repo_id)
  if !errors.Is(err, Error_Invalid) && err != nil {
    log.Printf("Error while removeing repo: %v\n", err.Error())
  }

  err = git.Delete_Repo(repo_id);
  if err != nil {
    log.Printf("Error while deleteing repo dir: %v\n", err.Error())
  }

  return err
}

func (d *Git_db) Get_Repos(
  ctx context.Context,
  user_id uint64,
  all bool,
  limit uint8,
  offset uint64,
) ([]string, error) {
  var max_limit uint8 = 100
  if limit > max_limit {
    return nil, Error_limit_too_big
  }

  repos, err := d.db.Get_Repos(ctx, user_id, all, limit, offset)
  if !errors.Is(err, Error_Invalid) && err != nil {
    log.Printf("Error while getting repos: %v\n", err.Error())
  }

  return repos, err
}

func (d *Git_db) Info(
  ctx context.Context,
  repo_id uint64,
) (*Repo_Info, error) {
  info, err := d.db.Info(ctx, repo_id)
  if !errors.Is(err, Error_Invalid) && err != nil {
    log.Printf("Error while getting repo info: %v\n", err.Error())
    return nil, err
  }

  return info, err
}
