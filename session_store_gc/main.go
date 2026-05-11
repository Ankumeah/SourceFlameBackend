package main

import (
	"context"
	"log"
	"os"
	"time"
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
  "SESSION_STORE_TYPE": "",
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
  keys := store.store.get_keys(Ctx)
  store.store.clean(Ctx, keys)
}

func connect_session_store(ctx context.Context) *store {
  switch env_vars["SESSION_STORE_TYPE"] {
    case "redis_cluster":
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

    case "redis_standalone":
      store, err := get_redis_driver(ctx,
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

    default:
      log.Fatalf("Unsupported session store type: %v\n", env_vars["SESSION_STORE_TYPE"])
      return nil
  }
}
