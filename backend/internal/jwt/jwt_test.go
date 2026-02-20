package jwt_test

import (
	"github.com/Ankumeah/DeltaBase/internal/jwt"

	"testing"
)

func Test_Issue_jwt(t *testing.T) {
  tests := []struct {
    username string
  } {
    { "beta" },
  }

  for _, test := range tests {
    jwt, err := jwt.Issue_jwt(test.username)
    if err != nil {
      t.Errorf("Error: %v\n", err.Error())
    } else if len(jwt) <= 0 {
      t.Errorf("JWT empty: %v\n", jwt)
    }

    t.Logf("Resulting JWT: %v\n", jwt)
  }
}

func Test_Validate_jwt(t *testing.T) {
  tests := []struct {
    username string
  } {
    { "beta" },
  }

  for _, test := range tests {
    token, err := jwt.Issue_jwt(test.username)
    if err != nil {
      t.Error("Error while issueing JWT")
    }

    res, err := jwt.Validate_jwt(test.username, token)
    if err != nil {
      t.Errorf("Error: %v\n", err.Error())
    } else if !res {
      t.Error("JWT was invalid")
    }
  }

}
