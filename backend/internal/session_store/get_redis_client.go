package session_store

import (
	"github.com/redis/go-redis/v9"

	"context"
)

func Get_Redis_Client(
  ctx context.Context,
  conn_config Universal_Redis_Config,
) (*redis.Client, error) {
  client := redis.NewClient(&redis.Options {
    Addr: conn_config.Hostname + ":" + conn_config.Port,
    Username: conn_config.Username,
    Password: conn_config.Password,
    DialTimeout: timeout,
    ReadTimeout: timeout,
    WriteTimeout: timeout,
  })
  if err := client.Ping(ctx).Err(); err != nil {
    return nil, err
  }

  return client, nil
}
