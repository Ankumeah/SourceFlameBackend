package session_store

import (
  "github.com/redis/go-redis/v9"

  "context"
  "time"
  "errors"
)

type redis_driver struct {
  rdb *redis.Client
}

func Get_Redis_Driver(ctx context.Context, url string) (*Session_store, error) {
  opts, err := redis.ParseURL(url)
  if err != nil {
    return &Session_store{}, errors.New("Invalid redis URL")
  }

  client := redis.NewClient(opts)
  _, err = client.Ping(ctx).Result()
  if err != nil {
    return &Session_store{}, err
  }

  driver := redis_driver { client }

  return &Session_store { &driver }, nil
}

func (r *redis_driver) Get(ctx context.Context, key string) (string, error) {
  res, err := r.rdb.Get(ctx, key).Result()
  if err == redis.Nil {
    return "", Error_not_found
  }

  return res, err
}

func (r *redis_driver) SetEx(ctx context.Context, key string, value string, expiration time.Duration) error {
  return r.rdb.SetEx(ctx, key, value, expiration).Err()
}

func (r *redis_driver) Del(ctx context.Context, keys ...string) (int64, error) {
  return r.rdb.Del(ctx, keys...).Result()
}
