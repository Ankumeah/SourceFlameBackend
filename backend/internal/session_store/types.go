package session_store

import (
	"context"
	"time"
)

type sessions_driver interface {
	Add_Session(ctx context.Context, username string, token string, timeout time.Duration) error
	Validate_Session(ctx context.Context, username string, token string) (bool, error)
	Delete_Session(ctx context.Context, username string, token string) error
	Get_Session_Count(ctx context.Context, username string) (int, error)
}

type Session_store struct {
	db              sessions_driver
	token_timeout   time.Duration
	token_limit     int
	token_length    int
	token_namespace string
}
