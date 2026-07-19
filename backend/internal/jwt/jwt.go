package jwt

import (
	"github.com/golang-jwt/jwt/v5"

	"log"
	"time"
)

func (h *JWTHandler) IssueJwt(username string) (string, error) {
	now := time.Now()

	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(h.lifespan)),
		},
	)

	signedJwt, err := t.SignedString(h.key)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return signedJwt, nil
}

func (h *JWTHandler) ValidateJwt(targetJwt string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(targetJwt, claims,
		func(t *jwt.Token) (any, error) { return h.key, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithJSONNumber(),
		jwt.WithIssuedAt(),
		jwt.WithLeeway(h.leeway),
	)
	if err != nil {
		return "", err
	} else if !token.Valid {
		return "", ErrInvalidJWT
	}

	return claims.Subject, nil
}
