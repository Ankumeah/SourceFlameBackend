package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"golang.org/x/crypto/argon2"
)

func Get_Hasher(time uint32, salt_len uint32, memory uint32, threads uint8, key_len uint32) *Hasher {
	return &Hasher{
		time:     time,
		salt_len: salt_len,
		memory:   memory,
		threads:  threads,
		key_len:  key_len,
	}
}

func (h *Hasher) Generate_Hash(password []byte, salt []byte) (*Hash, error) {
	var err error
	if len(salt) == 0 {
		salt, err = random_salt(h.salt_len)

		if err != nil {
			return nil, err
		}
	}

	_hash := argon2.IDKey(password, salt, h.time, h.memory, h.threads, h.key_len)

	hash := &Hash{Hash: _hash, Salt: salt}

	return hash, nil
}

func (h *Hasher) Compare_Hash(hash *Hash, password string) (bool, error) {
	_password_hash, err := h.Generate_Hash([]byte(password), hash.Salt)
	if err != nil {
		return false, err
	}

	password_hash := append(_password_hash.Hash, _password_hash.Salt...)
	target_hash := append(hash.Hash, hash.Salt...)

	if subtle.ConstantTimeCompare(target_hash, password_hash) != 1 {
		return false, nil
	} else {
		return true, nil
	}
}

func random_salt(length uint32) ([]byte, error) {
	secret := make([]byte, length)

	_, err := rand.Read(secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}
