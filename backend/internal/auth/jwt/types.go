package jwt

import (
  "errors"
)

type handler func(string) (string, error)

var Error_unsupported_JWT_type = errors.New("Unsupported JWT type")
var Error_invalid_JWT = errors.New("Invalid JWT")

var supported_JWT_types = map[string]handler {
  "google": google,
}
