package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"golang.org/x/crypto/argon2"
)

func GetHasher(time uint32, saltLen uint32, memory uint32, threads uint8, keyLen uint32) *Hasher {
	return &Hasher{
		time:    time,
		saltLen: saltLen,
		memory:  memory,
		threads: threads,
		keyLen:  keyLen,
	}
}

func (h *Hasher) GenerateHash(password []byte, salt []byte) (*Hash, error) {
	var err error
	if len(salt) == 0 {
		salt, err = randomSalt(h.saltLen)

		if err != nil {
			return nil, err
		}
	}

	_hash := argon2.IDKey(password, salt, h.time, h.memory, h.threads, h.keyLen)

	hash := &Hash{Hash: _hash, Salt: salt}

	return hash, nil
}

func (h *Hasher) CompareHash(hash *Hash, password string) (bool, error) {
	_passwordHash, err := h.GenerateHash([]byte(password), hash.Salt)
	if err != nil {
		return false, err
	}

	passwordHash := append(_passwordHash.Hash, _passwordHash.Salt...)
	targetHash := append(hash.Hash, hash.Salt...)

	if subtle.ConstantTimeCompare(targetHash, passwordHash) != 1 {
		return false, nil
	} else {
		return true, nil
	}
}

func randomSalt(length uint32) ([]byte, error) {
	secret := make([]byte, length)

	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}
