package database

import (
	"github.com/Ankumeah/SourceFlameBackend/internal/git"

	"context"
	"errors"
	"log"
	"sync"
)

func (d *Git_db) Create_Repo(
  ctx context.Context,
  owner_id uint64,
  repo_name string,
  private bool,
) (uint64, error) {
  repo_id, err := d.db.Create_Repo(ctx, owner_id, repo_name, private)
  if err != nil {
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
  var wg sync.WaitGroup

  var git_err error
  var db_err error

  wg.Go(func() {
    git_err = git.Delete_Repo(repo_id);
    if git_err != nil {
      log.Printf("Error while deleteing repo dir: %v\n", git_err.Error())
    }
  })
  wg.Go(func() {
    db_err = d.db.Delete_Repo(ctx, repo_id)
    if !errors.Is(db_err, Error_Invalid) && db_err != nil {
      log.Printf("Error while removeing repo: %v\n", db_err.Error())
    }
  })
  wg.Wait()

  if git_err != nil { return git_err }
  return db_err
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
