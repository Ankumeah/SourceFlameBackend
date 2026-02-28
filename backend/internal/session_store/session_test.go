package session_store_test

import (
	"github.com/Ankumeah/DeltaBase/internal/session_store"

	"context"
	"testing"
)

var Ctx = context.Background()

func Test_Add_Session(t *testing.T) {
  driver := session_store.Get_Fake_Driver()

  token, err := driver.Add_Session(Ctx, "test")
  if err != nil {
    t.Errorf("Error: %v\n", err.Error())
  }

  t.Logf("Resulting token: %v\n", token)
}

func Test_Validate_Session(t *testing.T) {
  driver := session_store.Get_Fake_Driver()

  token, err := driver.Add_Session(Ctx, "test")
  if err != nil {
    t.Errorf("Error: %v\n", err.Error())
  }

  valid, err := driver.Validate_Session(Ctx, "test", token)
  if err != nil {
    t.Errorf("Error: %v\n", err.Error())
  } else if !valid {
    t.Errorf("Invalid token")
  }

  t.Logf("Resulting token: %v\n", token)
}

func Test_Delete_Session(t *testing.T) {
  driver := session_store.Get_Fake_Driver()

  token, err := driver.Add_Session(Ctx, "test")
  if err != nil {
    t.Errorf("Error: %v\n", err.Error())
  }

  err = driver.Delete_Session(Ctx, "test", token)
  if err != nil {
    t.Errorf("Error: %v\n", err.Error())
  }
}
