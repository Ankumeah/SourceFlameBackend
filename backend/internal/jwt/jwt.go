package jwt

import (
	"github.com/golang-jwt/jwt/v5"

	"os"
	"time"
  "log"
  "strconv"
)

const jwt_leeway = 10 * time.Second

var jwt_key []byte
var jwt_lifespan time.Duration

func init() {
  key, ok := os.LookupEnv("JWT_KEY")
  if !ok { log.Fatalln("Unset env var: JWT_KEY") }
  jwt_key = []byte(key)

  lifespan, ok := os.LookupEnv("JWT_LIFESPAN_MINS")
  if !ok { log.Fatalln("Unset env var: JWT_LIFESPAN_MINS") }

  val, err := strconv.Atoi(lifespan)
  if err != nil { log.Fatalf("Error converting JWT_LIFESPAN_MINS to int: %v\n", err.Error()) }
  jwt_lifespan = time.Duration(val) * time.Minute
}

func Issue_jwt(username string) (string, error) {
  now := time.Now()

  t := jwt.NewWithClaims(jwt.SigningMethodHS256,
    jwt.RegisteredClaims {
      Subject: username,
      IssuedAt: jwt.NewNumericDate(now),
      ExpiresAt: jwt.NewNumericDate(now.Add(jwt_lifespan)),
    },
  )

  signed_jwt, err := t.SignedString(jwt_key)
  if err != nil {
    log.Println(err.Error())
    return "", err
  }

  return signed_jwt, nil
}

func Validate_jwt(target_jwt string) (string, error) {
  claims := &jwt.RegisteredClaims{}
  token, err := jwt.ParseWithClaims(target_jwt, claims,
    func(t *jwt.Token) (any, error) {
      return jwt_key, nil
    },
    jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
    jwt.WithJSONNumber(),
    jwt.WithIssuedAt(),
    jwt.WithLeeway(jwt_leeway),
  )
  if err != nil { return "", err }
  if !token.Valid { return "", Error_invalid_JWT }

  return claims.Subject, nil
}
