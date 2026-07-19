package session_store

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"sync"
)

func (d *SessionStore) AddSession(ctx context.Context, username string) (string, error) {
	var wg sync.WaitGroup

	var count int
	var countErr error
	var tokenErr error
	var token string

	wg.Go(func() {
		count, countErr = d.db.GetSessionCount(ctx, d.tokenNamespace+username)
		if countErr != nil {
			log.Printf("Error while getting session count: %v\n", countErr.Error())
		} else if count >= d.tokenLimit {
			countErr = ErrTooManyTokens
		}
	})
	wg.Go(func() {
		_token := make([]byte, d.tokenLength)
		if _, tokenErr = rand.Read(_token); tokenErr != nil {
			log.Printf("Error while generating token: %v\n", tokenErr.Error())
			return
		}
		token = base64.RawURLEncoding.EncodeToString(_token)
	})
	wg.Wait()
	if countErr != nil {
		return "", countErr
	}
	if tokenErr != nil {
		return "", tokenErr
	}

	if err := d.db.AddSession(
		ctx,
		d.tokenNamespace+username,
		token,
		d.tokenTimeout,
	); err != nil {
		log.Printf("Error while adding session: %v\n", err.Error())
		return "", err
	}

	return token, nil
}

func (d *SessionStore) ValidateSession(ctx context.Context, username string, token string) (bool, error) {
	valid, err := d.db.ValidateSession(ctx, d.tokenNamespace+username, token)
	if err != nil {
		log.Printf("Error while validating session: %v\n", err)
		return false, err
	}

	return valid, nil
}

func (d *SessionStore) DeleteSession(ctx context.Context, username string, token string) error {
	err := d.db.DeleteSession(ctx, d.tokenNamespace+username, token)
	if err != nil {
		log.Printf("Error while deleting session: %v\n", err.Error())
	}

	return err
}
