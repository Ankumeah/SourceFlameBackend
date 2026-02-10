package jwt

func Validate(JWT_type string, JWT string) (string, error) {
  validatation_func, ok := supported_JWT_types[JWT_type]
  if !ok {
    return "", Error_unsupported_JWT_type
  }

  return validatation_func(JWT)
}
