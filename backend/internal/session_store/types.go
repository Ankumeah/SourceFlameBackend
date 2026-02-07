package session_store

import (
  "context"
  "time"
  "errors"
)

type driver interface {
  Get(ctx context.Context, key string) (string, error)
  Del(ctx context.Context, keys ...string) (int64, error)
  SetEx(ctx context.Context, key string, value string, expiration time.Duration) error
}

var Error_not_found = errors.New("key not found")

type Session_store struct {
  db driver
}
