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
) ([]RepoInfo, error) {
	query := s.db.Rebind(`
    SELECT repos.name, repos.created_at, repos.stars, users.username
    FROM repos
    INNER JOIN users on repos.owner_id = users.id
    WHERE owner_id = ?
    LIMIT ? OFFSET ?
  `)

	rows, err := s.db.QueryContext(ctx, query, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repos []RepoInfo
	for rows.Next() {
		var repo RepoInfo
		err := rows.Scan(&repo.Name, &repo.Creation, &repo.Stars, &repo.Owner)
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
    SELECT repos.name, repos.created_at, repos.stars, users.username
    FROM repos
    INNER JOIN users ON repos.owner_id = users.id
    WHERE repos.id = ?;
  `)

	var info RepoInfo
	err := s.db.QueryRowContext(ctx, query, repoId).Scan(
		&info.Name, &info.Creation, &info.Stars, &info.Owner,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrInvalidRepo
	}

	return &info, err
}

func (s *gitSqlxDriver) TransferOwner(
	ctx context.Context,
	repoId uint64,
	newOwnerId uint64,
) error {
	query := s.db.Rebind(`
    UPDATE repos
    SET owner_id = ?
    WHERE id = ?;
  `)

	tag, err := s.db.ExecContext(ctx, query, newOwnerId, repoId)
	if err != nil {
		return err
	}

	affected, err := tag.RowsAffected()
	if affected == 0 {
		return ErrInvalidRepo
	}

	return nil
}

func (s *gitSqlxDriver) AddMember(
	ctx context.Context,
	repoId uint64,
	memberId uint64,
) error {
	query := s.db.Rebind(`
    INSERT INTO members (member_id, repo_id, addition)
    VALUES (?, ?, ?);
  `)

	now := time.Now().Unix()
	_, err := s.db.ExecContext(ctx, query, memberId, repoId, now)

	if isUniqueViolation(err) {
		err = ErrMemberExists
	}

	return err
}

func (s *gitSqlxDriver) RemoveMember(
	ctx context.Context,
	repoId uint64,
	memberId uint64,
) error {
	query := s.db.Rebind(`
    DELETE FROM members
    WHERE (repo_id = ? AND member_id = ?);
  `)

	tag, err := s.db.ExecContext(ctx, query, repoId, memberId)
	if err != nil {
		return err
	}

	affected, err := tag.RowsAffected()
	if affected == 0 {
		return ErrInvalidUser
	}

	return err
}

func (s *gitSqlxDriver) GetMembers(
	ctx context.Context,
	repoId uint64,
) ([]MemberInfo, error) {
	query := s.db.Rebind(`
    SELECT users.username, members.addition
    FROM members
    INNER JOIN users ON members.member_id = users.id
    WHERE members.repo_id = ?;
  `)

	rows, err := s.db.QueryContext(ctx, query, repoId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []MemberInfo
	for rows.Next() {
		var member MemberInfo
		err := rows.Scan(&member.Username, &member.Addition)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	if rows.Err() != nil {
		return nil, err
	}

	return members, err
}
