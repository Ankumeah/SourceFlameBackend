package database

import (
	"github.com/jmoiron/sqlx"

	"context"
	"database/sql"
	"time"
)

type git_sqlx_driver struct{ db *sqlx.DB }

func Git_Sql_Driver(db *sqlx.DB) *Git_db {
	return &Git_db{
		&git_sqlx_driver{db},
	}
}

func (s *git_sqlx_driver) Create_Repo(
	ctx context.Context,
	owner_id uint64,
	repo_name string,
	private bool,
) (uint64, error) {
	query := s.db.Rebind(`
    INSERT INTO repos (name, owner_id, private, created_at, stars)
    VALUES (?, ?, ?, ?, 0)
    RETURNING id;
  `)

	var repo_id uint64
	now := time.Now().Unix()
	err := s.db.QueryRowContext(ctx, query, repo_name, owner_id, private, now).Scan(&repo_id)

	if is_unique_violation(err) {
		err = Error_repo_exists
	}

	return repo_id, err
}

func (s *git_sqlx_driver) Get_Id(
	ctx context.Context,
	owner_id uint64,
	repo_name string,
) (uint64, error) {
	query := s.db.Rebind(`
    SELECT id FROM repos
    WHERE (owner_id = ? AND name = ?);
  `)

	var repo_id uint64
	err := s.db.QueryRowContext(ctx, query, owner_id, repo_name).Scan(&repo_id)
	if err == sql.ErrNoRows {
		return 0, Error_invalid_repo
	}

	return repo_id, err
}

func (s *git_sqlx_driver) Delete_Repo(
	ctx context.Context,
	repo_id uint64,
) error {
	query := s.db.Rebind(`
    DELETE FROM repos
    WHERE (id = ?);
  `)

	tag, err := s.db.ExecContext(ctx, query, repo_id)
	if err != nil {
		return err
	}

	effected, err := tag.RowsAffected()
	if effected == 0 {
		return Error_invalid_repo
	}

	return err
}

func (s *git_sqlx_driver) Get_Repos(
	ctx context.Context,
	user_id uint64,
	all bool,
	limit uint8,
	offset uint64,
) ([]string, error) {
	query_all := `
    SELECT name FROM repos
    WHERE (owner_id = ?)
    LIMIT ? OFFSET ?;
  `
	_query := `
    SELECT name FROM repos
    WHERE (NOT private AND owner_id = ?)
    LIMIT ? OFFSET ?
  `
	if all {
		_query = query_all
	}
	query := s.db.Rebind(_query)

	rows, err := s.db.QueryContext(ctx, query, user_id, limit, offset)
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

func (s *git_sqlx_driver) Info(
	ctx context.Context,
	repo_id uint64,
) (*Repo_Info, error) {
	query := s.db.Rebind(`
    SELECT repos.created_at, repos.stars, repos.private, users.username
    FROM repos
    INNER JOIN users ON repos.owner_id = users.id
    WHERE repos.id = ?;
  `)

	var info Repo_Info
	err := s.db.QueryRowContext(ctx, query, repo_id).Scan(
		&info.Creation, &info.Stars, &info.Private, &info.Owner,
	)
	if err == sql.ErrNoRows {
		return nil, Error_invalid_repo
	}

	return &info, err
}
