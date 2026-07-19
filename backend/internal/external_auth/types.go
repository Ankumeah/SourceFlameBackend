package external_auth

type handler func(string) (string, error)

var supportedJWTTypes = map[string]handler{
	"google": google,
}
