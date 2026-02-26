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

func (d *Session_store) Add_Session(ctx context.Context, username string) (string, error) {
  _token := make([]byte, token_length)
  if _, err := rand.Read(_token); err != nil {
    log.Printf("Error while genrating token: %v\n", err.Error())
    return "", err
  }
  token := base64.RawURLEncoding.EncodeToString(_token)

  if err := d.db.Add_Session(ctx, token_namespace + username, token, token_timeout); err != nil {
    log.Printf("Error while adding session: %v\n", err.Error())
    return "", err
  }

  return token, nil
}

func (d *Session_store) Validate_Session(ctx context.Context, username string, token string) (bool, error) {
  valid, err := d.db.Validate_Session(ctx, token_namespace + username, token)
  if err != nil {
    log.Printf("Error while validateing session: %v\n", err)
    return false, err
  }

  return valid, nil
}

func (d *Session_store) Delete_Session(ctx context.Context, username string, token string) error {
  if err := d.db.Delete_Session(ctx, username, token); err != nil {
    log.Printf("Error while deleting session: %v\n", err.Error())
    return err
  }

  return nil
}
