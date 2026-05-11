package main

import (
	a "github.com/Ankumeah/DeltaBase/internal/app"
	"github.com/Ankumeah/DeltaBase/internal/apis"
	"github.com/Ankumeah/DeltaBase/internal/database"
	"github.com/Ankumeah/DeltaBase/internal/middlewares"
	"github.com/Ankumeah/DeltaBase/internal/session_store"
	"github.com/Ankumeah/DeltaBase/internal/pat"

	"github.com/gin-gonic/gin"

	"context"
	"log"
	"os"
  "sync"
)

var Ctx = context.Background()
var app = &a.App{}
const pat_prefix = "gho_"
const pat_length = 32
var env_vars = map[string]string {
  "API_VERSION": "",
  "BACKEND_PORT": "",
  "SESSION_STORE_SESSIONS_USERNAME": "",
  "SESSION_STORE_SESSIONS_PASSWORD": "",
  "SESSION_STORE_HOSTNAME": "",
  "SESSION_STORE_PORT": "",
  "SESSION_STORE_TYPE": "",
  "DATABASE_USER": "",
  "DATABASE_PASSWORD": "",
  "DATABASE_HOST": "",
  "DATABASE_PORT": "",
  "DATABASE_DB": "",
  "DATABASE_CONFIG": "",
}

func load_env() {
  for env := range env_vars {
    _env, ok := os.LookupEnv(env)
    if !ok {
      log.Fatalf("Unset env var: %v\n", env)
    }

    env_vars[env] = _env
  }
}

func connect_session_store() {
  switch env_vars["SESSION_STORE_TYPE"] {
    case "redis_standalone":
      conn_config := session_store.Redis_Config {
        Username: env_vars["SESSION_STORE_SESSIONS_USERNAME"],
        Password: env_vars["SESSION_STORE_SESSIONS_PASSWORD"],
        Hostname: env_vars["SESSION_STORE_HOSTNAME"],
        Port: env_vars["SESSION_STORE_PORT"],
      }

      store, err := session_store.Get_Redis_Driver(Ctx, conn_config)
      if err != nil {
        log.Fatalf("Error while connecting to session store: %v\n", err.Error())
      }
      app.Store = store
    case "redis_cluster":
      conn_config := session_store.Redis_Cluster_Config {
        Username: env_vars["SESSION_STORE_SESSIONS_USERNAME"],
        Password: env_vars["SESSION_STORE_SESSIONS_PASSWORD"],
        Hostname: env_vars["SESSION_STORE_HOSTNAME"],
        Port: env_vars["SESSION_STORE_PORT"],
      }

      store, err := session_store.Get_Redis_Cluster_Driver(Ctx, conn_config)
      if err != nil {
        log.Fatalf("Error while connecting to session store: %v\n", err.Error())
      }
      app.Store = store
    default:
      log.Fatalf("Unsupported session store type: %v\n", env_vars["SESSION_STORE_TYPE"])
  }

  log.Println("Connected to session store")
}

func connect_database() {
  conn_config := database.Connection_Config {
    Username: env_vars["DATABASE_USER"],
    Password: env_vars["DATABASE_PASSWORD"],
    Hostname: env_vars["DATABASE_HOST"],
    Port: env_vars["DATABASE_PORT"],
    Db_name: env_vars["DATABASE_DB"],
    Db_config: env_vars["DATABASE_CONFIG"],
  }

  pool, err := database.Get_Connection_Pool(Ctx, conn_config)
  if err != nil {
    log.Fatalf("Error while connecting to database: %v\n", err.Error())
  }

  app.User_db = database.User_Postgres_Driver(pool)
  app.Git_db = database.Git_Postgres_Driver(pool)
  app.PAT_db = database.PAT_Postgres_Driver(pool)

  log.Println("Connected to database")
}

func get_pat_handler() {
  app.PAT_Handler = pat.Get_PAT_Handler(pat_prefix, pat_length)
  log.Println("Got PAT Handler")
}

func main() {
  log.Println("Loading env")
  load_env()

  var wg sync.WaitGroup
  wg.Go(func() {
	  log.Println("Connecting to session store")
    connect_session_store()
  })
  wg.Go(func() {
	  log.Println("Connecting to database")
    connect_database()
  })
  wg.Go(func() {
	  log.Println("Getting PAT Handler")
    get_pat_handler()
  })
  wg.Wait()

	log.Println("Starting http server")
	r := gin.Default()
	apiGroup := r.Group(
    "/api/" + env_vars["API_VERSION"] + "/",
    middlewars.Log_Middleware(),
  )
	apis.Apis(apiGroup, app)

  log.Println("Running backend on port: " + env_vars["BACKEND_PORT"])
  log.Println("API_VERSION: " + env_vars["API_VERSION"])
	r.Run(":" + env_vars["BACKEND_PORT"])
}
