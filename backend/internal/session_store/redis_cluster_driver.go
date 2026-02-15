package session_store

import (
	"github.com/redis/go-redis/v9"

	"context"
	"errors"
	"net"
	"time"
)

type redis_cluster_driver struct {
  rdb *redis.ClusterClient
}

func Get_Redis_Cluster_Driver(ctx context.Context, username string, password string, hostname string, port string) (*Session_store, error) {
  _addrs, err := net.LookupHost(hostname)
  if err != nil {
    return &Session_store{}, errors.New("Error while looking up ips: " + err.Error())
  }
  addrs := make([]string, len(_addrs))

  for i, addr := range _addrs {
    addrs[i] = addr + ":" + port
  }

  client := redis.NewClusterClient(&redis.ClusterOptions {
    Addrs: addrs,
    Username: username,
    Password: password,
  })
  _, err = client.Ping(ctx).Result()
  if err != nil {
    return &Session_store{}, err
  }

  driver := redis_cluster_driver { client }

  return &Session_store { &driver }, nil
}

func (r *redis_cluster_driver) Get(ctx context.Context, key string) (string, error) {
  res, err := r.rdb.Get(ctx, key).Result()
  if err == redis.Nil {
    return "", error_not_found
  }

  return res, err
}

func (r *redis_cluster_driver) SetEx(ctx context.Context, key string, value string, expiration time.Duration) error {
  return r.rdb.SetEx(ctx, key, value, expiration).Err()
}

func (r *redis_cluster_driver) Del(ctx context.Context, keys ...string) (int64, error) {
  return r.rdb.Del(ctx, keys...).Result()
}
