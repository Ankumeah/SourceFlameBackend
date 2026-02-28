package database

import (
	"github.com/Ankumeah/DeltaBase/internal/hash"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"context"
	"fmt"
)

const pool_max_conns = 10
const pool_min_conns = 2
const pool_min_idle_conns = 2

type postgres_driver struct {
  db *pgxpool.Pool
}

func Get_Postgres_Driver(
  ctx context.Context,
  username string,
  password string,
  hostname string,
  port string,
  dbname string,
) (*Database, error) {
  url_string := fmt.Sprintf(
    "user=%v password=%v host=%v port=%v dbname=%v pool_max_conns=%v pool_min_conns=%v pool_min_idle_conns=%v",
    username,
    password,
    hostname,
    port,
    dbname,
    pool_max_conns,
    pool_min_conns,
    pool_min_idle_conns,
  )
  config, err := pgxpool.ParseConfig(url_string)
  if err != nil {
    return nil, err
  }

  conn, err := pgxpool.NewWithConfig(ctx, config)
  if err != nil {
    return nil, err
  }

  driver := postgres_driver { conn }
  return &Database { &driver }, nil
}

func (p *postgres_driver) Add_User(ctx context.Context, username string, password_hash *hash.Hash) (uint64, error) {
  query := "INSERT INTO users (username, password_hash, salt) VALUES ($1, $2, $3) RETURNING (id);"

  user_id := new(uint64)
  if err := p.db.QueryRow(ctx, query, username, password_hash.Hash, password_hash.Salt).Scan(user_id); err != nil {
    return 0, err
  }

  return *user_id, nil
}

func (p *postgres_driver) Is_User_Valid(ctx context.Context, user_id uint64) (bool, error) {
  query := "SELECT (id) FROM users WHERE (id = $1);"

  err := p.db.QueryRow(ctx, query, user_id).Scan(new(uint64))
  if err == pgx.ErrNoRows {
    return false, nil
  } else if err != nil {
    return false, err
  } else {
    return true, nil
  }
}

func (p *postgres_driver) Delete_User(ctx context.Context, user_id uint64) error {
  query := "DELETE FROM users WHERE (id = $1);"

  if _, err := p.db.Exec(ctx, query, user_id); err != nil {
    return err
  }

  return nil
}

func (p *postgres_driver) Get_Hash(ctx context.Context, user_id uint64) (*hash.Hash, error) {
  query := "SELECT (password_hash, salt) FROM users WHERE (id = $1);"

  password_hash := new([]byte)
  salt := new([]byte)
  err := p.db.QueryRow(ctx, query, user_id).Scan(password_hash, salt)
  if err == pgx.ErrNoRows {
    return nil, Error_invalid_user
  } else if err != nil {
    return nil, err
  } else {
    return &hash.Hash {
      Hash: *password_hash,
      Salt: *salt,
    }, nil
  }
}

func (p *postgres_driver) Get_Id(ctx context.Context, username string) (uint64, error) {
  query := "SELECT (id) FROM users WHERE (username = $1);"

  id := new(uint64)
  err := p.db.QueryRow(ctx, query, username).Scan(id)
  if err == pgx.ErrNoRows {
    return 0, Error_invalid_user
  } else if err != nil {
    return 0, err
  } else {
    return *id, nil
  }
}
