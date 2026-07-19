package jwt_test

import (
	"github.com/Ankumeah/SourceFlameBackend/internal/jwt"

	"errors"
	"testing"
	"time"
)

func TestIssueJwt(t *testing.T) {
	tests := []struct {
		username string
	}{
		{"beta"},
	}
	handler := jwt.GetJWTHandler(time.Second, time.Minute, []byte("key"))

	for _, test := range tests {
		jwt, err := handler.IssueJwt(test.username)
		if err != nil {
			t.Errorf("Error: %v\n", err.Error())
		} else if len(jwt) <= 0 {
			t.Errorf("JWT empty: %v\n", jwt)
		}

		t.Logf("Resulting JWT: %v\n", jwt)
	}
}

func TestValidateJwt(t *testing.T) {
	tests := []struct {
		username string
	}{
		{"beta"},
	}
	handler := jwt.GetJWTHandler(time.Second, time.Minute, []byte("key"))

	for _, test := range tests {
		token, err := handler.IssueJwt(test.username)
		if err != nil {
			t.Error("Error while issuing JWT")
		}

		username, err := handler.ValidateJwt(token)
		if errors.Is(err, jwt.ErrInvalidJWT) {
			t.Errorf("JWT was invalid")
		} else if err != nil {
			t.Errorf("Error: %v\n", err.Error())
		} else if username != test.username {
			t.Error("username did not match")
		}
	}

}
