package middlewares

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

func parseBasicAuth(s string) (string, string, error) {
	_decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", "", err
	}

	decoded := string(_decoded)
	parts := strings.Split(decoded, ":")
	if len(parts) != 2 {
		return "", "", errors.New("Incorrect segments in auth header")
	}

	username := parts[0]
	password := parts[1]

	if username == "" || password == "" {
		return "", "", errors.New("Empty credentials")
	}

	return username, password, nil
}

func internalServerError() (int, map[string]any) {
	return http.StatusInternalServerError, map[string]any{"error": "Internal server error"}
}
func badHttpRequest(err error) (int, map[string]any) {
	return http.StatusBadRequest, map[string]any{"error": err.Error()}
}
func unauthorisedRequest() (int, map[string]any) {
	return http.StatusUnauthorized, map[string]any{"error": "Unauthorised"}
}
