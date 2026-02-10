package session_store

import (
  "context"
  "time"
)

type fake_driver struct {
  db map[string]string
}

func Get_Fake_Driver(ctx context.Context, url string) (*Session_store, error) {
  return &Session_store { &fake_driver{ map[string]string{} } }, nil
}

func (r *fake_driver) Get(ctx context.Context, key string) (string, error) {
  res, ok := r.db[key]
  if !ok  {
    return "", error_not_found
  }

  return res, nil
}

func (r *fake_driver) SetEx(ctx context.Context, key string, value string, expiration time.Duration) error {
  r.db[key] = value

  return nil
}

func (r *fake_driver) Del(ctx context.Context, keys ...string) (int64, error) {
  for _, key := range keys {
    r.db[key] = ""
  }

  return int64(len(keys)), nil
}
