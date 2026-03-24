package database

import (
  "github.com/Ankumeah/DeltaBase/internal/git"

  "context"
  "log"
)

func (d *Git_db) Create_Repo(ctx context.Context, username string, repo_name string, private bool) (uint64, error) {
  repo_id, err := d.db.Create_Repo(ctx, username, repo_name, private)
  if err != nil {
    log.Printf("Error while adding repo: %v\n", err.Error())
    return 0, err
  }

  err = git.Create_Repo(repo_id, private)
  if err != nil {
    log.Printf("Error while creating repo dir: %v\n", err.Error())
    d.db.Delete_Repo(ctx, repo_id)

    return 0, err
  }

  return repo_id, nil
}

func (d *Git_db) Get_Id(ctx context.Context, username string, repo_name string) (uint64, error) {
  repo_id, err := d.db.Get_Id(ctx, username, repo_name)
  if err != Error_invalid_repo && err != nil {
    log.Printf("Error while getting repo id: %v\n", err.Error())
    return 0, err
  }

  return repo_id, err
}

func (d *Git_db) Delete_Repo(ctx context.Context, repo_id uint64) error {
  err := d.db.Delete_Repo(ctx, repo_id)
  if err != Error_invalid_repo && err != nil {
    log.Printf("Error while deleteing repo id: %v\n", err.Error())
    return err
  }

  return nil
}
