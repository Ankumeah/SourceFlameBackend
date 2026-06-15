package database

import (
  "github.com/jmoiron/sqlx"

  "database/sql"
	"context"
	"errors"
	"time"
)

type pat_sqlx_driver struct { db *sqlx.DB }
func PAT_Sql_Driver(pool *sqlx.DB) *PAT_db {
  return &PAT_db {
    &pat_sqlx_driver { pool },
  }
}

func (s *pat_sqlx_driver) Add_PAT(ctx context.Context, owner_id uint64, hash string, pat_name string) (uint64, error) {
  query := s.db.Rebind(`
    INSERT INTO pats (name, hash, owner_id, created_at)
    VALUES (?, ?, ?, ?)
    RETURNING id;
  `)

  var pat_id uint64
  now := time.Now().Unix()
  err := s.db.QueryRowContext(ctx, query, pat_name, hash, owner_id, now).Scan(&pat_id)

  if is_unique_violation(err) { err = Error_pat_exists }

  return pat_id, err
}

func (s *pat_sqlx_driver) Validate_PAT(ctx context.Context, owner_id uint64, hash string) (uint64, error) {
  query := s.db.Rebind(`
    SELECT id FROM pats
    WHERE (owner_id = ? AND hash = ?)
  `)

  var pat_id uint64
  err := s.db.QueryRowContext(ctx, query, owner_id, hash).Scan(&pat_id)
  if errors.Is(err, sql.ErrNoRows) {
    err = Error_invalid_pat
  }

  return pat_id, err
}

func (s *pat_sqlx_driver) Get_Id(ctx context.Context, owner_id uint64, pat_name string) (uint64, error) {
  query := s.db.Rebind(`
    SELECT id FROM pats
    WHERE (owner_id = ? AND name = ?);
  `)

  var pat_id uint64
  err := s.db.QueryRowContext(ctx, query, owner_id, pat_name).Scan(&pat_id)
  if errors.Is(err, sql.ErrNoRows) {
    err = Error_invalid_pat
  }

  return pat_id, err
}

func (s *pat_sqlx_driver) Delete_PAT(ctx context.Context, pat_id uint64) error {
  query := s.db.Rebind(`
    DELETE FROM pats
    WHERE (id = ?);
  `)

  tag, err := s.db.ExecContext(ctx, query, pat_id)
  if err != nil { return err }

  effected, err := tag.RowsAffected()
  if effected == 0 { return Error_invalid_pat }

  return err
}
func (s *pat_sqlx_driver) Get_PATs(ctx context.Context, user_id uint64) ([]string, error) {
  query := s.db.Rebind(`
    SELECT name FROM pats
    WHERE (owner_id = ?);
  `)

  rows, err := s.db.QueryContext(ctx, query, user_id)
  if err != nil { return nil, err }
  defer rows.Close()

  var pats []string
  for rows.Next() {
    var pat string
    err := rows.Scan(&pat)
    if err != nil { return nil, err }
    pats = append(pats, pat)
  }
  if rows.Err() != nil { return nil, err }

  return pats, nil
}
func (s *pat_sqlx_driver) Info(ctx context.Context, pat_id uint64) (*PAT_Info, error) {
  query := s.db.Rebind(`
    SELECT name, created_at, last_used FROM pats
    WHERE (id = ?);
  `)

  var info PAT_Info
  err := s.db.QueryRowContext(ctx, query, pat_id).Scan(&info.Name, &info.Creation, &info.Last_Used)

  if err == sql.ErrNoRows {
    return nil, Error_invalid_pat
  } else if err != nil {
    return nil, err
  } else {
    return &info, nil
  }
}

func (s *pat_sqlx_driver) Update_Use(ctx context.Context, pat_id uint64) error {
  query := s.db.Rebind(`
    UPDATE pats
    SET last_used = ?
    WHERE id = ?
  `)
  now := time.Now().Unix()

  tag, err := s.db.ExecContext(ctx, query, now, pat_id)
  if err != nil { return err }

  effected, err := tag.RowsAffected()
  if err != nil {
    return err
  } else if effected <= 0 {
    return Error_invalid_pat
  } else { return nil }

}
