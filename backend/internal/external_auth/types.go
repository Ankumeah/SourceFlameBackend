package external_auth

type handler func(string) (string, error)

var supported_JWT_types = map[string]handler {
  "google": google,
}
