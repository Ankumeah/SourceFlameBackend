package session_store

import (
  "context"
  "time"
  "errors"
)

type driver interface {
  Add_Session(ctx context.Context, username string, token string, timeout time.Duration) error
  Validate_Session(ctx context.Context, username string, token string) (bool, error)
  Delete_Session(ctx context.Context, username string, token string) error
  Get_Session_Count(ctx context.Context, username string) (uint8, error)
}

type Session_store struct {
  db driver
}

var Error_too_many_tokens = errors.New("Too many tokens")
