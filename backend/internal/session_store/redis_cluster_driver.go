package session_store

import (
	"github.com/redis/go-redis/v9"

	"context"
	"errors"
	"net"
	"time"
)

const max_host_lookup_retires = 10
const host_lookup_rest = time.Second * 10

type redis_cluster_driver struct {
  rdb *redis.ClusterClient
}

func Get_Redis_Cluster_Driver(
  ctx context.Context,
  username string,
  password string,
  hostname string,
  port string,
) (*Session_store, error) {
  _addrs := []string {}
  err := errors.New("")
  for range max_host_lookup_retires {
    _addrs, err = net.LookupHost(hostname)
    if err == nil { break }
    time.Sleep(host_lookup_rest)
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
  })
  _, err = client.Ping(ctx).Result()
  if err != nil {
    return nil, err
  }

  driver := redis_cluster_driver { client }
  return &Session_store { &driver }, nil
}

func (r *redis_cluster_driver) Add_Session(ctx context.Context, token string, username string, timeout time.Duration) error {
  return r.rdb.SetEx(ctx, token_namespace + token, username, token_timeout).Err()
}

func (r *redis_cluster_driver) Validate_Session(ctx context.Context, username string, token string) (bool, error) {
  val, err := r.rdb.Get(ctx, token_namespace + token).Result()

  if err == redis.Nil {
    return false, nil
  } else if err != nil {
    return false, err
  } else if val != username {
    return false, nil
  } else {
    return true, nil
  }
}

func (r *redis_cluster_driver) Delete_Session(ctx context.Context, token string) error {
  _, err := r.rdb.Del(ctx, token_namespace + token).Result();
  return err
}
