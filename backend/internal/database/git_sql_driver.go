package database

import (
	"github.com/jmoiron/sqlx"

	"context"
	"database/sql"
	"errors"
	"time"
)

type gitSqlxDriver struct{ db *sqlx.DB }

func GitSqlDriver(db *sqlx.DB) *GitDb {
	return &GitDb{
		&gitSqlxDriver{db},
	}
}

func (s *gitSqlxDriver) CreateRepo(
	ctx context.Context,
	ownerId uint64,
	repoName string,
) (uint64, error) {
	query := s.db.Rebind(`
    INSERT INTO repos (name, owner_id, created_at, stars)
    VALUES (?, ?, ?, 0)
    RETURNING id;
  `)

	var repoId uint64
	now := time.Now().Unix()
	err := s.db.QueryRowContext(ctx, query, repoName, ownerId, now).Scan(&repoId)

	if isUniqueViolation(err) {
		err = ErrRepoExists
	}

	return repoId, err
}

func (s *gitSqlxDriver) GetId(
	ctx context.Context,
	ownerId uint64,
	repoName string,
) (uint64, error) {
	query := s.db.Rebind(`
    SELECT id FROM repos
    WHERE (owner_id = ? AND name = ?);
  `)

	var repoId uint64
	err := s.db.QueryRowContext(ctx, query, ownerId, repoName).Scan(&repoId)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrInvalidRepo
	}

	return repoId, err
}

func (s *gitSqlxDriver) DeleteRepo(
	ctx context.Context,
	repoId uint64,
) error {
	query := s.db.Rebind(`
    DELETE FROM repos
    WHERE (id = ?);
  `)

	tag, err := s.db.ExecContext(ctx, query, repoId)
	if err != nil {
		return err
	}

	affected, err := tag.RowsAffected()
	if affected == 0 {
		return ErrInvalidRepo
	}

	return err
}

func (s *gitSqlxDriver) GetRepos(
	ctx context.Context,
	userId uint64,
	limit uint8,
	offset uint64,
) ([]string, error) {
	query := s.db.Rebind(`
    SELECT name FROM repos
    WHERE (owner_id = ?)
    LIMIT ? OFFSET ?
  `)

	rows, err := s.db.QueryContext(ctx, query, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repos []string
	for rows.Next() {
		var repo string
		err := rows.Scan(&repo)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return repos, nil
}

func (s *gitSqlxDriver) Info(
	ctx context.Context,
	repoId uint64,
) (*RepoInfo, error) {
	query := s.db.Rebind(`
    SELECT repos.created_at, repos.stars, users.username
    FROM repos
    INNER JOIN users ON repos.owner_id = users.id
    WHERE repos.id = ?;
  `)

	var info RepoInfo
	err := s.db.QueryRowContext(ctx, query, repoId).Scan(
		&info.Creation, &info.Stars, &info.Owner,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrInvalidRepo
	}

	return &info, err
}
