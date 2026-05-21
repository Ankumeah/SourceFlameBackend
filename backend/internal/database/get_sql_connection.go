package database

import (
  "github.com/jmoiron/sqlx"

	"time"
	"context"
)

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

func Get_DB_Connection(ctx context.Context, url string, config Sql_Config) (*sqlx.DB, error) {
  db, err := sqlx.Open(driver_name, url)
  if err != nil { return nil, err }
  db.SetMaxOpenConns(config.max_conns)
  db.SetMaxIdleConns(config.max_idle)
  db.SetConnMaxLifetime(config.max_lifetime)
  db.SetConnMaxIdleTime(config.max_idle_time)

  if err := db.Ping(); err != nil { return nil, err }

  return db, nil
}
