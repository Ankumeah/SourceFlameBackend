package database

import (
	"github.com/Ankumeah/SourceFlameBackend/internal/hash"

  "github.com/jmoiron/sqlx"

  "database/sql"
	"context"
  "time"
)

type user_sqlx_driver struct { db *sqlx.DB }
func User_Sql_Driver(pool *sqlx.DB) *User_db {
  return &User_db {
    &user_sqlx_driver { pool },
  }
}

func (s *user_sqlx_driver) Add_User(
  ctx context.Context,
  username string,
  password_hash *hash.Hash,
) (uint64, error) {
  query := s.db.Rebind(`
    INSERT INTO users (username, created_at, password_hash, salt)
    VALUES (?, ?, ?, ?) RETURNING (id);
  `)

  var user_id uint64
  now := time.Now().Unix()
  if err := s.db.QueryRowContext(ctx, query,
    username, now, password_hash.Hash, password_hash.Salt,
  ).Scan(&user_id); err != nil {
    return 0, err
  }

  return user_id, nil
}

func (s *user_sqlx_driver) Delete_User(
  ctx context.Context,
  user_id uint64,
) error {
  query := s.db.Rebind("DELETE FROM users WHERE (id = ?);")

  _, err := s.db.ExecContext(ctx, query, user_id)
  return err
}

func (s *user_sqlx_driver) Get_Hash(
  ctx context.Context,
  user_id uint64,
) (*hash.Hash, error) {
  query := s.db.Rebind("SELECT password_hash, salt FROM users WHERE (id = ?);")

  var password_hash []byte
  var salt []byte
  err := s.db.QueryRowContext(ctx, query, user_id).Scan(&password_hash, &salt)
  if err == sql.ErrNoRows {
    return nil, Error_invalid_user
  } else if err != nil {
    return nil, err
  } else {
    return &hash.Hash {
      Hash: password_hash,
      Salt: salt,
    }, nil
  }
}

func (s *user_sqlx_driver) Get_Id(
  ctx context.Context,
  username string,
) (uint64, error) {
  query := s.db.Rebind("SELECT id FROM users WHERE (username = ?);")

  var id uint64
  err := s.db.QueryRowContext(ctx, query, username).Scan(&id)
  if err == sql.ErrNoRows {
    return 0, Error_invalid_user
  } else if err != nil {
    return 0, err
  } else {
    return id, nil
  }
}

func (s *user_sqlx_driver) Info(
  ctx context.Context,
  user_id uint64,
) (*User_Info, error) {
  query := s.db.Rebind("SELECT created_at FROM users WHERE (id = ?);")

  var creation uint64
  err := s.db.QueryRowContext(ctx, query, user_id).Scan(&creation)
  if err == sql.ErrNoRows {
    return nil, Error_invalid_user
  } else if err != nil {
    return nil, err
  } else {
    return &User_Info {
      Creation: creation,
    }, nil
  }
}
