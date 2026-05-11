package main

import "context"

type driver interface {
  get_keys(ctx context.Context) []string
  clean(ctx context.Context, keys []string) error
}
type store struct {
  store driver
}
