package session_store

import (
	"context"
	"time"
)

type fakeDriver struct {
	db map[string]string
}

func GetFakeDriver() *SessionStore {
	return &SessionStore{
		&fakeDriver{map[string]string{}},
		time.Minute,
		2, 8, "fake:",
	}
}

func (d *fakeDriver) AddSession(ctx context.Context, username string, token string, timeout time.Duration) error {
	d.db[token] = username
	return nil
}

func (d *fakeDriver) ValidateSession(ctx context.Context, username string, token string) (bool, error) {
	val, ok := d.db[token]

	if !ok {
		return false, nil
	} else if val != username {
		return false, nil
	} else {
		return true, nil
	}
}

func (d *fakeDriver) DeleteSession(ctx context.Context, username string, token string) error {
	return nil
}

func (d *fakeDriver) GetSessionCount(ctx context.Context, username string) (int, error) {
	return 1, nil
}
