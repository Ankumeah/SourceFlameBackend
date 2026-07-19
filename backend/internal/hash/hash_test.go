package hash_test

import (
	"github.com/Ankumeah/SourceFlameBackend/internal/hash"

	"testing"
)

func TestHash(t *testing.T) {
	tests := []struct {
		password []byte
		salt     []byte
	}{
		{[]byte("beta"), []byte("")},
	}
	var hasher = hash.GetHasher(2, 16, 64*1024, 4, 32)

	for _, test := range tests {
		hash, err := hasher.GenerateHash(test.password, test.salt)
		if err != nil {
			t.Errorf("Error: %v\n", err.Error())
		} else if len(hash.Hash) <= 0 {
			t.Error("hash empty")
		} else if len(hash.Salt) <= 0 {
			t.Error("salt empty")
		}

		t.Logf("Resulting hash: %v, salt: %v\n", string(hash.Hash), string(hash.Salt))

		valid, err := hasher.CompareHash(hash, string(test.password))
		if err != nil {
			t.Errorf("Error: %v\n", err.Error())
		} else if !valid {
			t.Errorf("hash not valid: %v, salt: %v\n", string(hash.Hash), string(hash.Salt))
		}
	}
}
