package session

import (
	"github.com/golang-jwt/jwt/v5"

	"os"
	"time"
  "log"
)

var jwt_key []byte

const jwt_lifespan = 15 * time.Minute
const jwt_leeway = 10 * time.Second

func init() {
  var key, ok = os.LookupEnv("JWT_KEY")
  if !ok {
    panic("env var JWT_KEY not set")
  }

  jwt_key = []byte(key)
}

func Issue_jwt(username string) (string, error) {
  now := time.Now()

  t := jwt.NewWithClaims(jwt.SigningMethodHS256,
    jwt.RegisteredClaims {
      Subject: username,
      IssuedAt: jwt.NewNumericDate(now),
      ExpiresAt: jwt.NewNumericDate(now.Add(jwt_lifespan)),
    })

  signed_jwt, err := t.SignedString(jwt_key)
  if err != nil {
    log.Println(err.Error())
    return "", err
  }

  return signed_jwt, nil
}

func Validate_jwt(username string, target_jwt string) (bool, error) {
  token, err := jwt.ParseWithClaims(target_jwt, &jwt.RegisteredClaims{},
    func(t *jwt.Token) (any, error) {
      return jwt_key, nil
    },
    jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
    jwt.WithSubject(username),
    jwt.WithJSONNumber(),
    jwt.WithIssuedAt(),
    jwt.WithLeeway(jwt_leeway),
  )

  if err != nil {
    return false, err
  }

  if !token.Valid {
    return  false, nil
  }

  return true, nil
}
