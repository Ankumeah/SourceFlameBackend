package jwt

import "time"

type JWT_Handler struct {
	leeway   time.Duration
	lifespan time.Duration
	key      []byte
}

func Get_JWT_Handler(
	leeway time.Duration,
	lifespan time.Duration,
	key []byte,
) *JWT_Handler {
	return &JWT_Handler{
		leeway:   leeway,
		lifespan: lifespan,
		key:      key,
	}
}
