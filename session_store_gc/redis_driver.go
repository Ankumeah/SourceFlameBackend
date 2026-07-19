package main

import (
	"github.com/redis/go-redis/v9"

  "context"
  "time"
  "strconv"
  "log"
)

type redis_driver struct {
  store *redis.Client
}
func get_redis_driver(
  ctx context.Context,
  username string,
  password string,
  hostname string,
  port string,
) (*store, error) {
  client := redis.NewClient(&redis.Options {
    Addr: hostname + ":" + port,
    Username: username,
    Password: password,
    DialTimeout: timeout,
    ReadTimeout: timeout,
    WriteTimeout: timeout,
  })
  if err := client.Ping(ctx).Err(); err != nil {
    return nil, err
  }

  driver := redis_driver { client }
  return &store { &driver }, nil
}

func (client *redis_driver) get_keys(ctx context.Context) []string {
  keys := []string {}

  _keys, cur, err := client.store.Scan(ctx, 0, namespace, 100).Result()
  if err != nil {
    log.Fatalf("Error while scanning: %v\n", err.Error())
  }
  keys = append(keys, _keys...)

  for cur != 0 {
    _keys, cur, err = client.store.Scan(ctx, cur, namespace, 100).Result()
    if err != nil {
      log.Fatalf("Error while scanning: %v\n", err.Error())
    }
    keys = append(keys, _keys...)
  }

  log.Println("Got keys")
  return keys
}

func (client *redis_driver) clean(ctx context.Context, keys []string) error {
  now := strconv.FormatInt(time.Now().Unix(), 10)

  for _, key := range keys {
    if _, err := client.store.ZRemRangeByScore(ctx, key, "-inf", now).Result(); err != nil {
      log.Fatalf("Error while cleaning: %v\n", err.Error())
    }
  }

  log.Println("Completed clean")
  return nil
}
