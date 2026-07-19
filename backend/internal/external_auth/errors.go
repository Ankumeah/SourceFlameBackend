package external_auth

import "errors"

var ErrUnsupportedJWTType = errors.New("Unsupported JWT type")
var ErrInvalidJWT = errors.New("Invalid JWT")
