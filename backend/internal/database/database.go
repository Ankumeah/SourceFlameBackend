package database

import (
  "github.com/Ankumeah/DeltaBase/internal/hash"

  "context"
  "log"
)

var hasher = hash.Get_Hasher(2, 16, 64 * 1024, 4, 32)

func (d *Database) Add_User(ctx context.Context, username string, password string) (uint64, error) {
  password_hash, err := hasher.Generate_Hash([]byte(password), []byte(""))
  if err != nil {
    log.Printf("Error while hashing password: %v\n", err.Error())
    return 0, err
  }

  user_id, err := d.db.Add_User(ctx, username, password_hash)
  if err != nil {
    log.Printf("Error while adding user: %v\n", err.Error())
  }

  return user_id, err
}

func (d *Database) Delete_User(ctx context.Context, user_id uint64) error {
  err := d.db.Delete_User(ctx, user_id)
  if err != nil {
    log.Printf("Error while deleting user: %v\n", err.Error())
  }

  return err
}

func (d *Database) Verify_User(ctx context.Context, user_id uint64, password string) (bool, error) {
  hash, err := d.db.Get_Hash(ctx, user_id)
  if err != Error_invalid_user && err != nil {
    log.Printf("Error while getting hash: %v\n", err.Error())
    return false, err
  }

  valid, err := hasher.Compare_Hash(hash, password)
  if err != nil {
    log.Printf("Error while compairing hash: %v\n", err.Error())
  }

  return valid, err
}

func (d *Database) Get_Id(ctx context.Context, username string) (uint64, error) {
  id, err := d.db.Get_Id(ctx, username)
  if err != Error_invalid_user && err != nil {
    log.Printf("Error while getting id: %v\n", err.Error())
  }

  return id, err
}
