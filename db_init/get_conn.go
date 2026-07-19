package main

import (
	"github.com/uptrace/bun"

	"context"
	"database/sql"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID           uint64 `bun:"id,pk,autoincrement"`
	Username     string `bun:"username,notnull,unique"`
	CreatedAt    uint64 `bun:"created_at,notnull"`
	PasswordHash []byte `bun:"password_hash,notnull"`
	Salt         []byte `bun:"salt,notnull"`
}
type Repo struct {
	bun.BaseModel `bun:"table:repos"`

	ID        uint64 `bun:"id,pk,autoincrement"`
	Name      string `bun:"name,notnull,unique:repo_name"`
	OwnerID   uint64 `bun:"owner_id,notnull,unique:repo_name"`
	CreatedAt uint64 `bun:"created_at,notnull"`
	Stars     uint64 `bun:"stars,notnull"`
}
type PAT struct {
	bun.BaseModel `bun:"table:pats"`

	ID        uint64 `bun:"id,pk,autoincrement"`
	Name      string `bun:"name,notnull,unique:pat_name"`
	Hash      []byte `bun:"hash,notnull"`
	OwnerID   uint64 `bun:"owner_id,notnull,unique:pat_name"`
	CreatedAt uint64 `bun:"created_at,notnull"`
	LastUsed  *int64 `bun:"last_used"`
}

type SqlConfig struct {
	maxConns    int
	maxIdle     int
	maxLifetime time.Duration
	maxIdleTime time.Duration
}

func NewSqlConfig(
	maxConn int,
	maxIdle int,
	maxLifetime time.Duration,
	maxIdleTime time.Duration,
) SqlConfig {
	return SqlConfig{
		maxConns:    maxConn,
		maxIdle:     maxIdle,
		maxLifetime: maxLifetime,
		maxIdleTime: maxIdleTime,
	}
}

func GetDBConnection(ctx context.Context, url string, config SqlConfig) (*sql.DB, error) {
	db, err := sql.Open(driverName, url)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.maxConns)
	db.SetMaxIdleConns(config.maxIdle)
	db.SetConnMaxLifetime(config.maxLifetime)
	db.SetConnMaxIdleTime(config.maxIdleTime)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func ExecSchema(ctx context.Context, db *sql.DB) error {
	bunDb := bun.NewDB(db, dia)

	tx, err := bunDb.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.NewCreateTable().Model((*User)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = tx.NewCreateTable().Model((*Repo)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = tx.NewCreateTable().Model((*PAT)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
