package database

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"context"
	"time"
)

type pat_postgres_driver struct {
  db *pgxpool.Pool
}

func PAT_Postgres_Driver( pool *pgxpool.Pool) *PAT_db {
  return &PAT_db {
    &pat_postgres_driver { pool },
  }
}

func (p *pat_postgres_driver) Create_PAT(ctx context.Context, owner_id uint64, hash string, pat_name string) (uint64, error) {
  query := `
    INSERT INTO pats (name, hash, owner_id, created_at)
    VALUES ($1, $2, $3, $4)
    RETURNING id;
  `

  var pat_id uint64
  now := time.Now().Unix()
  err := p.db.QueryRow(ctx, query, hash, owner_id, now).Scan(&pat_id)

  return pat_id, err
}

func (p *pat_postgres_driver) Get_Id(ctx context.Context, owner_id uint64, pat string) (uint64, error) {
  query := `
    SELECT id FROM pats
    WHERE (owner_id = $1 AND hash = $2);
  `

  var pat_id uint64
  err := p.db.QueryRow(ctx, query, owner_id, pat).Scan(&pat_id)

  return pat_id, err
}

func (p *pat_postgres_driver) Delete_PAT(ctx context.Context, pat_id uint64) error {
  query := `
    DELETE FROM pats
    WHERE (id = $1);
  `

  tag, err := p.db.Exec(ctx, query, pat_id)
  if tag.RowsAffected() == 0 {
    return Error_invalid_pat
  }

  return err
}
func (p *pat_postgres_driver) Get_PATs(ctx context.Context, user_id uint64) ([]string, error) {
  query := `
    SELECT name FROM pats
    WHERE (owner_id = $1);
  `

  rows, _ := p.db.Query(ctx, query, user_id)
  pats, err := pgx.CollectRows(rows, pgx.RowTo[string])

  return pats, err
}
func (p *pat_postgres_driver) Info(ctx context.Context, pat_id uint64) (*PAT_Info, error) {
  query := `
    SELECT name, created_at, last_used FROM pats
    WHERE (id = $1);
  `

  var info PAT_Info
  err := p.db.QueryRow(ctx, query, pat_id).Scan(&info.Name, &info.Creation, &info.Last_Used)

  return &info, err
}
