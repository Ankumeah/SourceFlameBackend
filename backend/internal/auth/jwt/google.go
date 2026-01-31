package jwt

import (
  "errors"
)

func Google(JWT string) (bool, error) {
  // TODO("Implement JWT varification")

  if JWT == "beta" {
    return true, nil
  }
  return false, errors.New("Not implemented")
}
