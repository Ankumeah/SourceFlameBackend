package middlewars

import (
	"encoding/base64"
	"strings"
)

func parse_basic_auth(s string) (string, string, error) {
	_decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil { return "", "", err }

	decoded := string(_decoded)
  parts := strings.Split(decoded, ":")
  if len(parts) != 2 {
    return "", "", invalid_auth_info
  }

  username := parts[0]
  pat := parts[1]

  if username == "" || pat  == "" {
    return "", "", empty_credintials
  }

  return username, pat, nil
}
