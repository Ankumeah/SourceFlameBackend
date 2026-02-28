package database

import (
  "github.com/Ankumeah/DeltaBase/internal/hash"

  "context"
  "errors"
)

type driver interface {
  Add_User(ctx context.Context, username string, password_hash *hash.Hash) (uint64, error)
  Get_Id(ctx context.Context, username string) (uint64, error)
  Is_User_Valid(ctx context.Context, user_id uint64) (bool, error)
  Delete_User(ctx context.Context, user_id uint64) error
  Get_Hash(ctx context.Context, user_id uint64) (*hash.Hash, error)
}

type Database struct {
  db driver
}

var Error_invalid_user = errors.New("Invalid user")
