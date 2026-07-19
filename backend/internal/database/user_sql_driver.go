package database

import (
	"github.com/Ankumeah/SourceFlameBackend/internal/hash"

	"github.com/jmoiron/sqlx"

	"context"
	"database/sql"
	"errors"
	"time"
)

type userSqlxDriver struct{ db *sqlx.DB }

func UserSqlDriver(pool *sqlx.DB) *UserDb {
	return &UserDb{
		&userSqlxDriver{pool},
	}
}

func (s *userSqlxDriver) AddUser(
	ctx context.Context,
	username string,
	passwordHash *hash.Hash,
) (uint64, error) {
	query := s.db.Rebind(`
    INSERT INTO users (username, created_at, password_hash, salt)
    VALUES (?, ?, ?, ?) RETURNING (id);
  `)

	var userId uint64
	now := time.Now().Unix()
	err := s.db.QueryRowContext(ctx, query,
		username, now, passwordHash.Hash, passwordHash.Salt,
	).Scan(&userId)

	if isUniqueViolation(err) {
		err = ErrUserExists
	}

	return userId, err
}

func (s *userSqlxDriver) DeleteUser(
	ctx context.Context,
	userId uint64,
) error {
	query := s.db.Rebind("DELETE FROM users WHERE (id = ?);")

	_, err := s.db.ExecContext(ctx, query, userId)
	return err
}

func (s *userSqlxDriver) GetHash(
	ctx context.Context,
	userId uint64,
) (*hash.Hash, error) {
	query := s.db.Rebind("SELECT password_hash, salt FROM users WHERE (id = ?);")

	var passwordHash []byte
	var salt []byte
	err := s.db.QueryRowContext(ctx, query, userId).Scan(&passwordHash, &salt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrInvalidUser
	} else if err != nil {
		return nil, err
	} else {
		return &hash.Hash{
			Hash: passwordHash,
			Salt: salt,
		}, nil
	}
}

func (s *userSqlxDriver) GetId(
	ctx context.Context,
	username string,
) (uint64, error) {
	query := s.db.Rebind("SELECT id FROM users WHERE (username = ?);")

	var id uint64
	err := s.db.QueryRowContext(ctx, query, username).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrInvalidUser
	} else if err != nil {
		return 0, err
	} else {
		return id, nil
	}
}

func (s *userSqlxDriver) Info(
	ctx context.Context,
	userId uint64,
) (*UserInfo, error) {
	query := s.db.Rebind("SELECT created_at FROM users WHERE (id = ?);")

	var creation uint64
	err := s.db.QueryRowContext(ctx, query, userId).Scan(&creation)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrInvalidUser
	} else if err != nil {
		return nil, err
	} else {
		return &UserInfo{
			Creation: creation,
		}, nil
	}
}
