package jwt

import "time"

type JWTHandler struct {
	leeway   time.Duration
	lifespan time.Duration
	key      []byte
}

func GetJWTHandler(
	leeway time.Duration,
	lifespan time.Duration,
	key []byte,
) *JWTHandler {
	return &JWTHandler{
		leeway:   leeway,
		lifespan: lifespan,
		key:      key,
	}
}
