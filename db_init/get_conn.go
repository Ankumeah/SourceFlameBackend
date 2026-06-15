package main

import (
	"github.com/uptrace/bun"

	"time"
	"context"
  "database/sql"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID uint64 `bun:"id,pk,autoincrement"`
	Username string `bun:"username,notnull,unique"`
	Created_At uint64 `bun:"created_at,notnull"`
	Password_Hash []byte `bun:"password_hash,notnull"`
	Salt []byte `bun:"salt,notnull"`
}
type Repo struct {
	bun.BaseModel `bun:"table:repos"`

	ID uint64 `bun:"id,pk,autoincrement"`
  Name string `bun:"name,notnull,unique:repo_name"`
	Private bool `bun:"private,notnull"`
	OwnerID uint64 `bun:"owner_id,notnull,unique:repo_name"`
	CreatedAt uint64 `bun:"created_at,notnull"`
	Stars uint64 `bun:"stars,notnull"`
}
type PAT struct {
	bun.BaseModel `bun:"table:pats"`

	ID uint64 `bun:"id,pk,autoincrement"`
  Name string `bun:"name,notnull,unique:pat_name"`
	Hash []byte `bun:"hash,notnull"`
	OwnerID uint64 `bun:"owner_id,notnull,unique:pat_name"`
	CreatedAt uint64 `bun:"created_at,notnull"`
	LastUsed *int64 `bun:"last_used"`
}

type Sql_Config struct {
  max_conns int
  max_idle int
  max_lifetime time.Duration
  max_idle_time time.Duration
}
func New_Sql_Config(
  max_conn int,
  max_idle int,
  max_lifetime time.Duration,
  max_idle_time time.Duration,
) Sql_Config {
  return Sql_Config {
    max_conns: max_conn,
    max_idle: max_idle,
    max_lifetime: max_lifetime,
    max_idle_time: max_idle_time,
  }
}

func Get_DB_Connection(ctx context.Context, url string, config Sql_Config) (*sql.DB, error) {
  db, err := sql.Open(driver_name, url)
  if err != nil { return nil, err }
  db.SetMaxOpenConns(config.max_conns)
  db.SetMaxIdleConns(config.max_idle)
  db.SetConnMaxLifetime(config.max_lifetime)
  db.SetConnMaxIdleTime(config.max_idle_time)

  if err := db.Ping(); err != nil { return nil, err }

  return db, nil
}

func Exec_Schema(ctx context.Context, db *sql.DB) error {
  bun_db := bun.NewDB(db, dia)

  tx, err := bun_db.BeginTx(ctx, &sql.TxOptions{})
  if err != nil { return err }
  defer tx.Rollback()

  _, err = tx.NewCreateTable().Model((*User)(nil)).IfNotExists().Exec(ctx)
  if err != nil { return err }
  _, err = tx.NewCreateTable().Model((*Repo)(nil)).IfNotExists().Exec(ctx)
  if err != nil { return err }
  _, err = tx.NewCreateTable().Model((*PAT)(nil)).IfNotExists().Exec(ctx)
  if err != nil { return err }

  return tx.Commit()
}
