package database

import (
	"github.com/Ankumeah/DeltaBase/internal/hash"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"context"
	"fmt"
  "time"
)

type user_postgres_driver struct {
  db *pgxpool.Pool
}

func User_Postgres_Driver(
  ctx context.Context,
  username string,
  password string,
  hostname string,
  port string,
  dbname string,
  db_config string,
) (*User_db, error) {
  url_string := fmt.Sprintf(
    "user=%v password=%v host=%v port=%v dbname=%v %v",
    username,
    password,
    hostname,
    port,
    dbname,
    db_config,
  )
  config, err := pgxpool.ParseConfig(url_string)
  if err != nil {
    return nil, err
  }

  conn, err := pgxpool.NewWithConfig(ctx, config)
  if err != nil {
    return nil, err
  }

  driver := user_postgres_driver { conn }
  return &User_db { &driver }, nil
}

func (p *user_postgres_driver) Add_User(
  ctx context.Context,
  username string,
  password_hash *hash.Hash,
) (uint64, error) {
  query := `
    INSERT INTO users (username, created_at, password_hash, salt)
    VALUES ($1, $2, $3, $4) RETURNING (id);
  `

  var user_id uint64
  now := time.Now().Unix()
  if err := p.db.QueryRow(ctx, query,
    username, now, password_hash.Hash, password_hash.Salt,
  ).Scan(&user_id); err != nil {
    return 0, err
  }

  return user_id, nil
}

func (p *user_postgres_driver) Is_User_Valid(
  ctx context.Context,
  user_id uint64,
) (bool, error) {
  query := "SELECT (id) FROM users WHERE (id = $1);"

  err := p.db.QueryRow(ctx, query, user_id).Scan()
  if err == pgx.ErrNoRows {
    return false, nil
  } else if err != nil {
    return false, err
  } else {
    return true, nil
  }
}

func (p *user_postgres_driver) Delete_User(
  ctx context.Context,
  user_id uint64,
) error {
  query := "DELETE FROM users WHERE (id = $1);"

  _, err := p.db.Exec(ctx, query, user_id)
  return err
}

func (p *user_postgres_driver) Get_Hash(
  ctx context.Context,
  user_id uint64,
) (*hash.Hash, error) {
  query := "SELECT (password_hash, salt) FROM users WHERE (id = $1);"

  var password_hash []byte
  var salt []byte
  err := p.db.QueryRow(ctx, query, user_id).Scan(&password_hash, &salt)
  if err == pgx.ErrNoRows {
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

func (p *user_postgres_driver) Get_Id(
  ctx context.Context,
  username string,
) (uint64, error) {
  query := "SELECT (id) FROM users WHERE (username = $1);"

  var id uint64
  err := p.db.QueryRow(ctx, query, username).Scan(&id)
  if err == pgx.ErrNoRows {
    return 0, Error_invalid_user
  } else if err != nil {
    return 0, err
  } else {
    return id, nil
  }
}

func (p *user_postgres_driver) Info(
  ctx context.Context,
  user_id uint64,
) (*User_Info, error) {
  query := "SELECT (created_at) FROM users WHERE (id = $1);"

  var creation uint64
  err := p.db.QueryRow(ctx, query, user_id).Scan(&creation)
  if err == pgx.ErrNoRows {
    return nil, Error_invalid_user
  } else if err != nil {
    return nil, err
  } else {
    return &User_Info {
      Creation: creation,
    }, nil
  }
}
