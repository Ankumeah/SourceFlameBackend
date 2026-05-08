package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5"

	"context"
  "time"
)

type git_postgres_driver struct {
  db *pgxpool.Pool
}

func Git_Postgres_Driver( pool *pgxpool.Pool) *Git_db {
  return &Git_db {
    &git_postgres_driver { pool },
  }
}

func (p *git_postgres_driver) Create_Repo(
  ctx context.Context,
  owner_id uint64,
  repo_name string,
  private bool,
) (uint64, error) {
  query := `
    INSERT INTO repos (name, owner_id, private, created_at, stars)
    VALUES ($1, $2, $3, $4, 0)
    RETURNING id;
  `

  var repo_id uint64
  now := time.Now().Unix()
  err := p.db.QueryRow(ctx, query, repo_name, owner_id, private, now).Scan(&repo_id)

  return repo_id, err
}

func (p *git_postgres_driver) Get_Id(
  ctx context.Context,
  owner_id uint64,
  repo_name string,
) (uint64, error) {
  query := `
    SELECT id FROM repos
    WHERE (owner_id = $1 AND name = $2);
  `

  var repo_id uint64
  err := p.db.QueryRow(ctx, query, owner_id, repo_name).Scan(&repo_id)
  if err == pgx.ErrNoRows {
    return 0, Error_invalid_repo
  }

  return repo_id, err
}

func (p *git_postgres_driver) Delete_Repo(
  ctx context.Context,
  repo_id uint64,
) error {
  query := `
    DELETE FROM repos
    WHERE (id = $1);
  `

  tag, err := p.db.Exec(ctx, query, repo_id)
  if tag.RowsAffected() == 0 {
    return Error_invalid_repo
  }
  return err
}

func (p *git_postgres_driver) Get_Repos(
  ctx context.Context,
  user_id uint64,
  all bool,
  limit uint8,
  offset uint64,
) ([]string, error) {
  query_all := `
    SELECT name FROM repos
    WHERE (owner_id = $1)
    LIMIT $2 OFFSET $3;
  `
  query := `
    SELECT name FROM repos
    WHERE (NOT private AND owner_id = $1)
    LIMIT $2 OFFSET $3
  `
  if all { query = query_all }

  rows, _ := p.db.Query(ctx, query, user_id, limit, offset)
  repos, err := pgx.CollectRows(rows, pgx.RowTo[string])
  if err != nil {
    return nil, err
  }

  return repos, nil
}

func (p *git_postgres_driver) Info(
  ctx context.Context,
  repo_id uint64,
) (*Repo_Info, error) {
  query := `
    SELECT repos.created_at, repos.stars, repos.private, users.username
    FROM repos
    INNER JOIN users ON repos.owner_id = users.id
    WHERE repos.id = $1;
  `

  var info Repo_Info
  err := p.db.QueryRow(ctx, query, repo_id).Scan(
    &info.Creation, &info.Stars, &info.Private, &info.Owner,
  )
  if err == pgx.ErrNoRows {
    return nil, Error_invalid_repo
  }

  return &info, err
}
