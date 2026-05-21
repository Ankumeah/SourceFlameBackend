package jwt

import (
	"github.com/golang-jwt/jwt/v5"

	"time"
  "log"
)

func (h *JWT_Handler) Issue_jwt(username string) (string, error) {
  now := time.Now()

  t := jwt.NewWithClaims(jwt.SigningMethodHS256,
    jwt.RegisteredClaims {
      Subject: username,
      IssuedAt: jwt.NewNumericDate(now),
      ExpiresAt: jwt.NewNumericDate(now.Add(h.lifespan)),
    },
  )

  signed_jwt, err := t.SignedString(h.key)
  if err != nil {
    log.Println(err.Error())
    return "", err
  }

  return signed_jwt, nil
}

func (h *JWT_Handler) Validate_jwt(target_jwt string) (string, error) {
  claims := &jwt.RegisteredClaims{}
  token, err := jwt.ParseWithClaims(target_jwt, claims,
    func(t *jwt.Token) (any, error) { return h.key, nil },
    jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
    jwt.WithJSONNumber(),
    jwt.WithIssuedAt(),
    jwt.WithLeeway(h.leeway),
  )
  if err != nil {
    return "", err
  } else if !token.Valid {
    return "", Error_invalid_JWT
  }

  return claims.Subject, nil
}
