package main

import (
	"github.com/redis/go-redis/v9"

  "context"
  "strconv"
  "time"
  "net"
  "errors"
  "log"
)

type redis_cluster_driver struct {
  store *redis.ClusterClient
}
func get_redis_cluster_driver(
  ctx context.Context,
  username string,
  password string,
  hostname string,
  port string,
) (*store, error) {
  _addrs := []string {}
  var err error
  for range max_host_lookup_retries {
    _addrs, err = net.LookupHost(hostname)
    if err == nil { break }
    time.Sleep(timeout)
  }

  if err != nil {
    return nil, errors.New("Error while looking up ips: " + err.Error())
  }

  addrs := make([]string, len(_addrs))
  for i, addr := range _addrs {
    addrs[i] = addr + ":" + port
  }

  client := redis.NewClusterClient(&redis.ClusterOptions {
    Addrs: addrs,
    Username: username,
    Password: password,
    DialTimeout: timeout,
    ReadTimeout: timeout,
    WriteTimeout: timeout,
  })

  if err = client.Ping(ctx).Err(); err != nil {
    return nil, err
  }

  driver := redis_cluster_driver { client }
  return &store { &driver }, nil
}

func (client *redis_cluster_driver) get_keys(ctx context.Context) []string {
  keys := []string {}

  if err := client.store.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) error {
    _keys, cur, err := client.Scan(ctx, 0, namespace, 100).Result()
    if err != nil { return err }
    keys = append(keys, _keys...)

    for cur != 0 {
      _keys, cur, err = client.Scan(ctx, cur, namespace, 100).Result()
      if err != nil { return err }
      keys = append(keys, _keys...)
    }

    return nil
  }); err != nil {
    log.Fatalf("Error while scanning: %v\n", err.Error())
  }

  log.Println("Got keys")
  return keys
}

func (client *redis_cluster_driver) clean(ctx context.Context, keys []string) error {
  now := strconv.FormatInt(time.Now().Unix(), 10)

  for _, key := range keys {
    if _, err := client.store.ZRemRangeByScore(ctx, key, "-inf", now).Result(); err != nil {
      log.Fatalf("Error while cleaning: %v\n", err.Error())
    }
  }

  log.Println("Completed clean")
  return nil
}
