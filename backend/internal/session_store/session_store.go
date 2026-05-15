package session_store

import (
  "context"
  "crypto/rand"
  "encoding/base64"
  "strconv"
  "time"
  "log"
  "os"
  "sync"
)

const token_length = 64

var token_timeout time.Duration
var token_namespace string
var token_limit uint8

func init() {
  timeout, ok := os.LookupEnv("TOKEN_TIMEOUT_DAYS")
  if !ok { log.Fatalln("Unset env var: TOKEN_TIMEOUT_DAYS") }

  val, err := strconv.Atoi(timeout)
  if err != nil { log.Fatalf("Error while converting TOKEN_TIMEOUT_DAYS to int: %v\n", err.Error()) }
  token_timeout = time.Duration(val) * 24 * time.Hour

  limit, ok := os.LookupEnv("TOKEN_LIMIT")
  if !ok { log.Fatalf("Unset env var: TOKEN_LIMIT") }

  val, err = strconv.Atoi(limit)
  if err != nil { log.Fatalf("Error while converting TOKEN_LIMIT to int: %v\n", err.Error()) }
  token_limit = uint8(val)

  namespace, ok := os.LookupEnv("TOKEN_NAMESPACE")
  if !ok { log.Fatalln("Unset env var: TOKEN_NAMESPACE") }
  token_namespace = namespace
}

func (d *Session_store) Add_Session(ctx context.Context, username string) (string, error) {
  var wg sync.WaitGroup

  var count uint8
  var count_err error
  var token_err error
  var token string

  wg.Go(func() {
    count, count_err = d.db.Get_Session_Count(ctx, token_namespace + username)
    if count_err != nil {
      log.Printf("Error while getting session count: %v\n", count_err.Error())
    } else if count >= token_limit {
      count_err = Error_too_many_tokens
    }
  })
  wg.Go(func() {
    _token := make([]byte, token_length)
    if _, token_err = rand.Read(_token); token_err != nil {
      log.Printf("Error while genrating token: %v\n", token_err.Error())
      return
    }
    token = base64.RawURLEncoding.EncodeToString(_token)
  })
  wg.Wait()
  if count_err != nil { return "", count_err }
  if token_err != nil { return "", token_err }

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
  if err := d.db.Delete_Session(ctx, token_namespace + username, token); err != nil {
    log.Printf("Error while deleting session: %v\n", err.Error())
  }

  return nil
}
