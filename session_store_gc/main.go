package main

import (
  "github.com/redis/go-redis/v9"

  "strconv"
  "context"
  "errors"
  "time"
  "log"
  "net"
  "os"
)

const max_host_lookup_retires = 10
const timeout = 10 * time.Second
const namespace = "refresh:"

var Ctx = context.Background()

var env_vars = map[string]string {
  "SESSION_STORE_GC_USERNAME": "",
  "SESSION_STORE_GC_PASSWORD": "",
  "SESSION_STORE_HOSTNAME": "",
  "SESSION_STORE_PORT": "",
}

func init() {
  for env := range env_vars {
    _env, ok := os.LookupEnv(env)
    if !ok {
      log.Fatalf("Unset env var: %v\n", env)
    }

    env_vars[env] = _env
  }
}

func main() {
  store := connect_session_store(Ctx)
  keys := get_keys(Ctx, store)
  clean(Ctx, store, keys)
}

func connect_session_store(ctx context.Context) *redis.ClusterClient {
  store, err := get_redis_cluster_driver(ctx,
    env_vars["SESSION_STORE_GC_USERNAME"],
    env_vars["SESSION_STORE_GC_PASSWORD"],
    env_vars["SESSION_STORE_HOSTNAME"],
    env_vars["SESSION_STORE_PORT"],
  )
  if err != nil {
    log.Fatalf("Error while connecting to session store: %v\n", err.Error())
  }

  log.Println("Connected to session store")
  return store
}

func get_redis_cluster_driver(
  ctx context.Context,
  username string,
  password string,
  hostname string,
  port string,
) (*redis.ClusterClient, error) {
  _addrs := []string {}
  var err error
  for range max_host_lookup_retires {
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

  return client, nil
}

func get_keys(ctx context.Context, client *redis.ClusterClient) []string {
  keys := []string {}

  if err := client.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) error {
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
    log.Fatalf("Error while scaning: %v\n", err.Error())
  }

  log.Println("Got keys")
  return keys
}

func clean(ctx context.Context, client *redis.ClusterClient, keys []string) error {
  now := strconv.FormatInt(time.Now().Unix(), 10)

  for _, key := range keys {
    if _, err := client.ZRemRangeByScore(ctx, key, "-inf", now).Result(); err != nil {
      log.Fatalf("Error while cleaning: %v\n", err.Error())
    }
  }

  log.Println("Complete clean")
  return nil
}
