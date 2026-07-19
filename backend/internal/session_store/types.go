package session_store

import (
	"context"
	"time"
)

type sessionsDriver interface {
	AddSession(ctx context.Context, username string, token string, timeout time.Duration) error
	ValidateSession(ctx context.Context, username string, token string) (bool, error)
	DeleteSession(ctx context.Context, username string, token string) error
	GetSessionCount(ctx context.Context, username string) (int, error)
}

type SessionStore struct {
	db             sessionsDriver
	tokenTimeout   time.Duration
	tokenLimit     int
	tokenLength    int
	tokenNamespace string
}
