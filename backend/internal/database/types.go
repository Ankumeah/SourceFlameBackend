package database

import (
  "github.com/Ankumeah/DeltaBase/internal/hash"

  "context"
  "errors"
)

type user_driver interface {
  Add_User(ctx context.Context, username string, password_hash *hash.Hash) (uint64, error)
  Get_Id(ctx context.Context, username string) (uint64, error)
  Is_User_Valid(ctx context.Context, user_id uint64) (bool, error)
  Delete_User(ctx context.Context, user_id uint64) error
  Get_Hash(ctx context.Context, user_id uint64) (*hash.Hash, error)
  Get_Creation(ctx context.Context, user_id uint64) (uint64, error)
}

type git_driver interface {
  Create_Repo(ctx context.Context, username string, repo_name string, private bool) (uint64, error)
  Get_Id(ctx context.Context, username string, repo_name string) (uint64, error)
  Delete_Repo(ctx context.Context, repo_id uint64) error
  Get_Creation(ctx context.Context, repo_id uint64) (uint64, error)
}

type User_db struct { db user_driver }
type Git_db struct { db git_driver }

var Error_invalid_user = errors.New("Invalid user")
var Error_invalid_repo = errors.New("Invalid repo")
