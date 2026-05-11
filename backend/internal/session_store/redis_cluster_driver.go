package session_store

import (
	"github.com/redis/go-redis/v9"

	"context"
	"errors"
	"net"
	"time"
  "strconv"
)

const max_host_lookup_retires = 10
const timeout = time.Second * 10

type redis_cluster_driver struct {
  rdb *redis.ClusterClient
}
type Redis_Cluster_Config struct {
  Username string
  Password string
  Hostname string
  Port string
}

func Get_Redis_Cluster_Driver(
  ctx context.Context,
  conn_config Redis_Cluster_Config,
) (*Session_store, error) {
  _addrs := []string {}
  var err error
  for range max_host_lookup_retires {
    _addrs, err = net.LookupHost(conn_config.Hostname)
    if err == nil { break }
    time.Sleep(timeout)
  }

  if err != nil {
    return nil, errors.New("Error while looking up ips: " + err.Error())
  }

  addrs := make([]string, len(_addrs))
  for i, addr := range _addrs {
    addrs[i] = addr + ":" + conn_config.Port
  }

  client := redis.NewClusterClient(&redis.ClusterOptions {
    Addrs: addrs,
    Username: conn_config.Username,
    Password: conn_config.Password,
    DialTimeout: timeout,
    ReadTimeout: timeout,
    WriteTimeout: timeout,
  })

  if err = client.Ping(ctx).Err(); err != nil {
    return nil, err
  }

  driver := redis_cluster_driver { client }
  return &Session_store { &driver }, nil
}

func (r *redis_cluster_driver) Add_Session(ctx context.Context, username string, token string, timeout time.Duration) error {
  now := time.Now().Add(timeout).Unix()
  return r.rdb.ZAdd(ctx, username, redis.Z { Member: token, Score: float64(now) }).Err()
}

func (r *redis_cluster_driver) Validate_Session(ctx context.Context, username string, token string) (bool, error) {
  now := float64(time.Now().Unix())
  exp, err := r.rdb.ZScore(ctx, username, token).Result()

  if err == redis.Nil {
    return false, nil
  } else if err != nil {
    return false, err
  } else if exp <= now {
    return false, nil
  } else {
    return true, nil
  }
}

func (r *redis_cluster_driver) Delete_Session(ctx context.Context, username string, token string) error {
  return r.rdb.ZRem(ctx, username, token).Err()
}

func (r *redis_cluster_driver) Get_Session_Count(ctx context.Context, username string) (uint8, error) {
  now := time.Now().Unix()
  count , err := r.rdb.ZCount(ctx, username, strconv.FormatInt(now, 10), "inf").Result()
  if err != nil {
    return 0, err
  } else {
    return uint8(count), nil
  }
}
