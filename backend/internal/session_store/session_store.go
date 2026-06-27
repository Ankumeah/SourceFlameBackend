package session_store

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"sync"
)

func (d *Session_store) Add_Session(ctx context.Context, username string) (string, error) {
	var wg sync.WaitGroup

	var count int
	var count_err error
	var token_err error
	var token string

	wg.Go(func() {
		count, count_err = d.db.Get_Session_Count(ctx, d.token_namespace+username)
		if count_err != nil {
			log.Printf("Error while getting session count: %v\n", count_err.Error())
		} else if count >= d.token_limit {
			count_err = Error_too_many_tokens
		}
	})
	wg.Go(func() {
		_token := make([]byte, d.token_length)
		if _, token_err = rand.Read(_token); token_err != nil {
			log.Printf("Error while genrating token: %v\n", token_err.Error())
			return
		}
		token = base64.RawURLEncoding.EncodeToString(_token)
	})
	wg.Wait()
	if count_err != nil {
		return "", count_err
	}
	if token_err != nil {
		return "", token_err
	}

	if err := d.db.Add_Session(
		ctx,
		d.token_namespace+username,
		token,
		d.token_timeout,
	); err != nil {
		log.Printf("Error while adding session: %v\n", err.Error())
		return "", err
	}

	return token, nil
}

func (d *Session_store) Validate_Session(ctx context.Context, username string, token string) (bool, error) {
	valid, err := d.db.Validate_Session(ctx, d.token_namespace+username, token)
	if err != nil {
		log.Printf("Error while validateing session: %v\n", err)
		return false, err
	}

	return valid, nil
}

func (d *Session_store) Delete_Session(ctx context.Context, username string, token string) error {
	err := d.db.Delete_Session(ctx, d.token_namespace+username, token)
	if err != nil {
		log.Printf("Error while deleting session: %v\n", err.Error())
	}

	return err
}
