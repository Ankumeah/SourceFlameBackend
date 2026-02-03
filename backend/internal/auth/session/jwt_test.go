package session_test

import (
	"github.com/Ankumeah/DeltaBase/internal/auth/session"

	"testing"
)

func Test_Issue_jwt(t *testing.T) {
  tests := []struct {
    user_id uint64
  } {
    {1},
  }

  for _, test := range tests {
    jwt, err := session.Issue_jwt(test.user_id)
    if err != nil {
      t.Errorf("Error: %v", err.Error())
    } else if len(jwt) <= 0 {
      t.Errorf("JWT empty: %v", jwt)
    }

    t.Logf("Resulting JWT: %v", jwt)
  }
}

func Test_Validate_jwt(t *testing.T) {
  tests := []struct {
    user_id uint64
  } {
    {1},
  }

  for _, test := range tests {
    jwt, err := session.Issue_jwt(test.user_id)
    if err != nil {
      t.Error("Error while issueing JWT")
    }

    res, err := session.Validate_jwt(test.user_id, jwt)
    if err != nil {
      t.Errorf("Error: %v", err.Error())
    } else if !res {
      t.Error("JWT was invalid")
    }
  }

}













