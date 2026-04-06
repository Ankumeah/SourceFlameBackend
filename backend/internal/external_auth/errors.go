package external_auth

import "errors"

var Error_unsupported_JWT_type = errors.New("Unsupported JWT type")
var Error_invalid_JWT = errors.New("Invalid JWT")
