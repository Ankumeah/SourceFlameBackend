package session_store_test

import (
	"github.com/Ankumeah/SourceFlameBackend/internal/session_store"

	"context"
	"testing"
)

var Ctx = context.Background()

func TestAddSession(t *testing.T) {
	driver := session_store.GetFakeDriver()

	token, err := driver.AddSession(Ctx, "test")
	if err != nil {
		t.Errorf("Error: %v\n", err.Error())
	}

	t.Logf("Resulting token: %v\n", token)
}

func TestValidateSession(t *testing.T) {
	driver := session_store.GetFakeDriver()

	token, err := driver.AddSession(Ctx, "test")
	if err != nil {
		t.Errorf("Error: %v\n", err.Error())
	}

	valid, err := driver.ValidateSession(Ctx, "test", token)
	if err != nil {
		t.Errorf("Error: %v\n", err.Error())
	} else if !valid {
		t.Errorf("Invalid token")
	}

	t.Logf("Resulting token: %v\n", token)
}

func TestDeleteSession(t *testing.T) {
	driver := session_store.GetFakeDriver()

	token, err := driver.AddSession(Ctx, "test")
	if err != nil {
		t.Errorf("Error: %v\n", err.Error())
	}

	err = driver.DeleteSession(Ctx, "test", token)
	if err != nil {
		t.Errorf("Error: %v\n", err.Error())
	}
}
