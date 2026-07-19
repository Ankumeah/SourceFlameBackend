package database

import (
	"github.com/jmoiron/sqlx"

	"context"
	"database/sql"
	"errors"
	"time"
)

type patSqlxDriver struct{ db *sqlx.DB }

func PATSqlDriver(pool *sqlx.DB) *PATDb {
	return &PATDb{
		&patSqlxDriver{pool},
	}
}

func (s *patSqlxDriver) AddPAT(ctx context.Context, ownerId uint64, hash string, patName string) (uint64, error) {
	query := s.db.Rebind(`
    INSERT INTO pats (name, hash, owner_id, created_at)
    VALUES (?, ?, ?, ?)
    RETURNING id;
  `)

	var patId uint64
	now := time.Now().Unix()
	err := s.db.QueryRowContext(ctx, query, patName, hash, ownerId, now).Scan(&patId)

	if isUniqueViolation(err) {
		err = ErrPatExists
	}

	return patId, err
}

func (s *patSqlxDriver) ValidatePAT(ctx context.Context, ownerId uint64, hash string) (uint64, error) {
	query := s.db.Rebind(`
    SELECT id FROM pats
    WHERE (owner_id = ? AND hash = ?)
  `)

	var patId uint64
	err := s.db.QueryRowContext(ctx, query, ownerId, hash).Scan(&patId)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrInvalidPat
	}

	return patId, err
}

func (s *patSqlxDriver) GetId(ctx context.Context, ownerId uint64, patName string) (uint64, error) {
	query := s.db.Rebind(`
    SELECT id FROM pats
    WHERE (owner_id = ? AND name = ?);
  `)

	var patId uint64
	err := s.db.QueryRowContext(ctx, query, ownerId, patName).Scan(&patId)
	if errors.Is(err, sql.ErrNoRows) {
		err = ErrInvalidPat
	}

	return patId, err
}

func (s *patSqlxDriver) DeletePAT(ctx context.Context, patId uint64) error {
	query := s.db.Rebind(`
    DELETE FROM pats
    WHERE (id = ?);
  `)

	tag, err := s.db.ExecContext(ctx, query, patId)
	if err != nil {
		return err
	}

	affected, err := tag.RowsAffected()
	if affected == 0 {
		return ErrInvalidPat
	}

	return err
}
func (s *patSqlxDriver) GetPATs(ctx context.Context, userId uint64) ([]string, error) {
	query := s.db.Rebind(`
    SELECT name FROM pats
    WHERE (owner_id = ?);
  `)

	rows, err := s.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pats []string
	for rows.Next() {
		var pat string
		err := rows.Scan(&pat)
		if err != nil {
			return nil, err
		}
		pats = append(pats, pat)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return pats, nil
}
func (s *patSqlxDriver) Info(ctx context.Context, patId uint64) (*PATInfo, error) {
	query := s.db.Rebind(`
    SELECT name, created_at, last_used FROM pats
    WHERE (id = ?);
  `)

	var info PATInfo
	err := s.db.QueryRowContext(ctx, query, patId).Scan(&info.Name, &info.Creation, &info.LastUsed)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrInvalidPat
	} else if err != nil {
		return nil, err
	} else {
		return &info, nil
	}
}

func (s *patSqlxDriver) UpdateUse(ctx context.Context, patId uint64) error {
	query := s.db.Rebind(`
    UPDATE pats
    SET last_used = ?
    WHERE id = ?
  `)
	now := time.Now().Unix()

	tag, err := s.db.ExecContext(ctx, query, now, patId)
	if err != nil {
		return err
	}

	affected, err := tag.RowsAffected()
	if err != nil {
		return err
	} else if affected <= 0 {
		return ErrInvalidPat
	} else {
		return nil
	}

}
