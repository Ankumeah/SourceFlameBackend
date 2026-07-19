package external_auth

func Validate(JWTType string, JWT string) (string, error) {
	validationFunc, ok := supportedJWTTypes[JWTType]
	if !ok {
		return "", ErrUnsupportedJWTType
	}

	return validationFunc(JWT)
}
