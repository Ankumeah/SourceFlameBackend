package database

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"context"
	"fmt"
  "log"
)

func Get_Connection_Pool(
  ctx context.Context,
  conn_config Connection_Config,
) (*pgxpool.Pool, error) {
  url_string := fmt.Sprintf(
    "user=%v password=%v host=%v port=%v dbname=%v %v",
    conn_config.Username,
    conn_config.Password,
    conn_config.Hostname,
    conn_config.Port,
    conn_config.Db_name,
    conn_config.Db_config,
  )
  config, err := pgxpool.ParseConfig(url_string)
  if err != nil {
    log.Printf("Error while parseing DB config: %v\n", err.Error())
    return nil, err
  }

  conn, err := pgxpool.NewWithConfig(ctx, config)
  if err != nil {
    log.Printf("Error while connecting to DB config: %v\n", err.Error())
    return nil, err
  }

  err = conn.Ping(ctx)
  if err != nil {
    return nil, err
  }

  return conn, nil
}
