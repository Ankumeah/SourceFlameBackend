package session

import (
	"github.com/golang-jwt/jwt/v5"

	"os"
	"time"
  "log"
  "strconv"
)

var _jwt_key, ok = os.LookupEnv("JWT_KEY")
var jwt_key = []byte(_jwt_key)

var jwt_lifespan = 15 * time.Minute
var jwt_leeway = 10 * time.Second

func init() {
  if !ok {
    panic("env var JWT_KEY not set")
  }
}

func Issue_jwt(user_id uint64) (string, error) {
  now := time.Now()

  t := jwt.NewWithClaims(jwt.SigningMethodHS256,
    jwt.RegisteredClaims {
      Subject: strconv.FormatUint(user_id, 10),
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

func Validate_jwt(user_id uint64, target_jwt string) (bool, error) {
  token, err := jwt.ParseWithClaims(target_jwt, &jwt.RegisteredClaims{},
    func(t *jwt.Token) (any, error) {
      return jwt_key, nil
    },
    jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
    jwt.WithSubject(strconv.FormatUint(user_id, 10)),
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
