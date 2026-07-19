package database

import (
	"github.com/Ankumeah/SourceFlameBackend/internal/hash"

	"context"
	"errors"
	"log"
)

var hasher = hash.GetHasher(2, 16, 64*1024, 4, 32)

func (d *UserDb) AddUser(
	ctx context.Context,
	username string,
	password string,
) (uint64, error) {
	passwordHash, err := hasher.GenerateHash([]byte(password), []byte(""))
	if err != nil {
		log.Printf("Error while hashing password: %v\n", err.Error())
		return 0, err
	}

	userId, err := d.db.AddUser(ctx, username, passwordHash)
	if !errors.Is(err, ErrSafe) && err != nil {
		log.Printf("Error while adding user: %v\n", err.Error())
	}

	return userId, err
}

func (d *UserDb) DeleteUser(ctx context.Context, userId uint64) error {
	err := d.db.DeleteUser(ctx, userId)
	if !errors.Is(err, ErrInvalid) && err != nil {
		log.Printf("Error while deleting user: %v\n", err.Error())
	}

	return err
}

func (d *UserDb) VerifyUser(
	ctx context.Context,
	userId uint64,
	password string,
) (bool, error) {
	hash, err := d.db.GetHash(ctx, userId)
	if !errors.Is(err, ErrInvalid) && err != nil {
		log.Printf("Error while getting hash: %v\n", err.Error())
		return false, err
	}

	valid, err := hasher.CompareHash(hash, password)
	if err != nil {
		log.Printf("Error while comparing hash: %v\n", err.Error())
	}

	return valid, err
}

func (d *UserDb) GetId(
	ctx context.Context,
	username string,
) (uint64, error) {
	id, err := d.db.GetId(ctx, username)
	if !errors.Is(err, ErrInvalid) && err != nil {
		log.Printf("Error while getting user id: %v\n", err.Error())
	}

	return id, err
}

func (d *UserDb) Info(
	ctx context.Context,
	userId uint64,
) (*UserInfo, error) {
	info, err := d.db.Info(ctx, userId)
	if !errors.Is(err, ErrInvalid) && err != nil {
		log.Printf("Error while getting user creation time: %v\n", err.Error())
	}

	return info, err
}
