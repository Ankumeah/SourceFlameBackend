package session_store

import (
  "context"
  "crypto/rand"
  "encoding/base64"
  "time"
  "log"
)

const token_length = 64
const token_timeout = 24 * 30 * time.Hour
const token_namespace = "refresh:"

func (d *Session_store) Add_Session(ctx context.Context, email string) (string, error) {
  _token := make([]byte, token_length)
  if _, err := rand.Read(_token); err != nil {
    log.Printf("Error while genrating token: %v\n", err.Error())
    return "", err
  }

  token := base64.RawURLEncoding.EncodeToString(_token)

  if err := d.db.SetEx(ctx, token_namespace + token, email, token_timeout); err != nil {
    log.Printf("Error while adding session: %v\n", err.Error())
    return "", err
  }

  return token, nil
}

func (d *Session_store) Validate_Session(ctx context.Context, email string, token string) (bool, error) {
  val, err := d.db.Get(ctx, token_namespace + token)
  if err == error_not_found {
    return false, nil
  } else if err != nil {
    log.Printf("Error while validating session: %v\n", err.Error())
    return false, err
  }

  if val != email {
    return false, nil
  }

  return true, nil
}

func (d *Session_store) Delete_Session(ctx context.Context, token string) error {
  if _, err := d.db.Del(ctx, token_namespace + token); err != nil {
    log.Printf("Error while deleting session: %v\n", err.Error())
  }

  return nil
}
