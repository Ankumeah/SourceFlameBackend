package pat_test

import (
	"github.com/Ankumeah/DeltaBase/internal/pat"

	"testing"
)

func Test_Genrate_PAT(t *testing.T) {
  handler := pat.Get_PAT_Handler("test_", 32)

  pat, err := handler.Genrate_PAT()
  if err != nil {
    t.Errorf("Error: %v\n", err.Error())
  }

  t.Logf("Resulting PAT: %v", string(pat))
}
