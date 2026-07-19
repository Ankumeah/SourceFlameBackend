package pat_test

import (
	"github.com/Ankumeah/SourceFlameBackend/internal/pat"

	"testing"
)

func TestGeneratePAT(t *testing.T) {
	handler := pat.GetPATHandler("test_", 32)

	pat, err := handler.GeneratePAT()
	if err != nil {
		t.Errorf("Error: %v\n", err.Error())
	}

	t.Logf("Resulting PAT: %v", string(pat))
}
