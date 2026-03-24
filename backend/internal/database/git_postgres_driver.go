package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"

	"context"
	"fmt"
  "time"
)

type git_postgres_driver struct {
  db *pgxpool.Pool
}

func Git_Postgres_Driver(
  ctx context.Context,
  username string,
  password string,
  hostname string,
  port string,
  dbname string,
  db_config string,
) (*Git_db, error) {
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

  driver := git_postgres_driver { conn }
  return &Git_db { &driver }, nil
}

func (p *git_postgres_driver) Create_Repo(ctx context.Context, username string, repo_name string, private bool) (uint64, error) {
  query := "INSERT INTO repos (name, owner, private, created_at, stars) VALUES ($1, $2, $3, $4, 0) RETURNING id;"

  var repo_id uint64
  now := time.Now().Unix()
  err := p.db.QueryRow(ctx, query, repo_name, username, private, now).Scan(&repo_id)

  return repo_id, err
}

func (p *git_postgres_driver) Get_Id(ctx context.Context, username string, repo_name string) (uint64, error) {
  query := "SELECT (id) FROM repos WHERE (owner = $1 AND name = $2);"

  var repo_id uint64
  err := p.db.QueryRow(ctx, query, username, repo_name).Scan(&repo_id)
  if err == pgx.ErrNoRows {
    return 0, Error_invalid_repo
  }

  return repo_id, err
}

func (p *git_postgres_driver) Delete_Repo(ctx context.Context, repo_id uint64) error {
  query := "DELETE FROM repos WHERE (id = $1);"

  _, err := p.db.Exec(ctx, query, repo_id)
  return err
}

func (p *git_postgres_driver) Get_Creation(ctx context.Context, repo_id uint64) (uint64, error) {
  query := "SELECT (created_at) FROM repos WHERE (id = $1);"

  var created_at uint64
  err := p.db.QueryRow(ctx, query, repo_id).Scan(&created_at)
  if err == pgx.ErrNoRows {
    return 0, Error_invalid_repo
  } else if err != nil {
    return 0, err
  } else {
    return created_at, nil
  }
}
