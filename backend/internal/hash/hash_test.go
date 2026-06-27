package hash_test

import (
	"github.com/Ankumeah/SourceFlameBackend/internal/hash"

	"testing"
)

func Test_Hash(t *testing.T) {
	tests := []struct {
		password []byte
		salt     []byte
	}{
		{[]byte("beta"), []byte("")},
	}
	var hasher = hash.Get_Hasher(2, 16, 64*1024, 4, 32)

	for _, test := range tests {
		hash, err := hasher.Generate_Hash(test.password, test.salt)
		if err != nil {
			t.Errorf("Error: %v\n", err.Error())
		} else if len(hash.Hash) <= 0 {
			t.Error("hash empty")
		} else if len(hash.Salt) <= 0 {
			t.Error("salt empty")
		}

		t.Logf("Resulting hash: %v, salt: %v\n", string(hash.Hash), string(hash.Salt))

		valid, err := hasher.Compare_Hash(hash, string(test.password))
		if err != nil {
			t.Errorf("Error: %v\n", err.Error())
		} else if !valid {
			t.Errorf("hash not valid: %v, salt: %v\n", string(hash.Hash), string(hash.Salt))
		}
	}
}
