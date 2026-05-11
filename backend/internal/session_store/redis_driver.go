package session_store

import (
	"github.com/redis/go-redis/v9"

	"context"
	"time"
  "strconv"
)

type redis_driver struct {
  rdb *redis.Client
}
type Redis_Config struct {
  Username string
  Password string
  Hostname string
  Port string
}

func Get_Redis_Driver(
  ctx context.Context,
  conn_config Redis_Config,
) (*Session_store, error) {
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

  driver := redis_driver { client }
  return &Session_store { &driver }, nil
}

func (r *redis_driver) Add_Session(ctx context.Context, username string, token string, timeout time.Duration) error {
  now := time.Now().Add(timeout).Unix()
  return r.rdb.ZAdd(ctx, username, redis.Z { Member: token, Score: float64(now) }).Err()
}

func (r *redis_driver) Validate_Session(ctx context.Context, username string, token string) (bool, error) {
  now := float64(time.Now().Unix())
  exp, err := r.rdb.ZScore(ctx, username, token).Result()

  if err == redis.Nil {
    return false, nil
  } else if err != nil {
    return false, err
  } else if exp <= now {
    return false, nil
  } else {
    return true, nil
  }
}

func (r *redis_driver) Delete_Session(ctx context.Context, username string, token string) error {
  return r.rdb.ZRem(ctx, username, token).Err()
}

func (r *redis_driver) Get_Session_Count(ctx context.Context, username string) (uint8, error) {
  now := time.Now().Unix()
  count , err := r.rdb.ZCount(ctx, username, strconv.FormatInt(now, 10), "inf").Result()
  if err != nil {
    return 0, err
  } else {
    return uint8(count), nil
  }
}
