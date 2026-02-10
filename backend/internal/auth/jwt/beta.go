package jwt

import (
  "errors"
)

func beta(JWT string) (string, error) {
  // TODO("Remove in prod")

  return "", errors.New("This is the test handler")
}
