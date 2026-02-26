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

func (d *fake_driver) Add_Session(ctx context.Context, token string, username string, timeout time.Duration) error {
  d.db[token] = username
  return nil
}

func (d *fake_driver) Validate_Session(ctx context.Context, username string, token string) (bool, error) {
  val, ok := d.db[token]

  if !ok {
    return false, nil
  } else if val != username {
    return false, nil
  } else {
    return true, nil
  }
}

func (d *fake_driver) Delete_Session(ctx context.Context, username string, token string) error {
  return nil
}

func (d *fake_driver) Get_Session_Count(ctx context.Context, username string) (uint8, error) {
  return 1, nil
}
