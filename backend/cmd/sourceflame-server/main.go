package main

import (
	"strconv"
	"time"

	"github.com/Ankumeah/SourceFlameBackend/internal/apis"
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/database"
	"github.com/Ankumeah/SourceFlameBackend/internal/middlewares"
	"github.com/Ankumeah/SourceFlameBackend/internal/pat"
	"github.com/Ankumeah/SourceFlameBackend/internal/session_store"

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
  "DATABASE_URL": "",
  "DATABASE_MAX_CONNS": "",
  "DATABASE_MAX_IDLE_CONNS": "",
  "DATABASE_MAX_LIFETIME": "",
  "DATABASE_MAX_IDLE_TIME": "",
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
  sessions_config := session_store.Universal_Redis_Config {
    Username: env_vars["SESSION_STORE_SESSIONS_USERNAME"],
    Password: env_vars["SESSION_STORE_SESSIONS_PASSWORD"],
    Hostname: env_vars["SESSION_STORE_HOSTNAME"],
    Port: env_vars["SESSION_STORE_PORT"],
  }
  var client session_store.Redis_client
  var err error

  switch env_vars["SESSION_STORE_TYPE"] {
    case "redis_standalone":
      client, err = session_store.Get_Redis_Client(Ctx, sessions_config)
    case "redis_cluster":
      client, err = session_store.Get_Redis_Cluster_Client(Ctx, sessions_config)
    default:
      log.Fatalf("Unsupported session store type: %v\n", env_vars["SESSION_STORE_TYPE"])
  }
  if err != nil {
    log.Panicf("Error while connecting to session store: %v\n", err.Error())
  }

  app.Store = session_store.Get_Sessions_Uinversal_Redis_Driver(client)
  log.Println("Connected to session store")
}

func connect_database() {
  max_conn, err := strconv.Atoi(env_vars["DATABASE_MAX_CONNS"])
  if err != nil {
    log.Fatalf("Error while parsing DATABASE_MAX_CONNS: %v\n", err.Error())
  }
  max_idle, err := strconv.Atoi(env_vars["DATABASE_MAX_IDLE_CONNS"])
  if err != nil {
    log.Fatalf("Error while parsing DATABASE_MAX_IDLE_CONNS: %v\n", err.Error())
  }
  max_lifetime, err := time.ParseDuration(env_vars["DATABASE_MAX_LIFETIME"])
  if err != nil {
    log.Fatalf("Error while parsing DATABASE_MAX_LIFETIME: %v\n", err.Error())
  }
  max_idle_time, err := time.ParseDuration(env_vars["DATABASE_MAX_IDLE_TIME"])
  if err != nil {
    log.Fatalf("Error while parsing DATABASE_MAX_IDLE_TIME: %v\n", err.Error())
  }
  config := database.New_Sql_Config(int(max_conn), int(max_idle), max_lifetime, max_idle_time)

  db, err := database.Get_DB_Connection(Ctx, env_vars["DATABASE_URL"], config)
  if err != nil {
    log.Fatalf("Error while connecting to database: %v\n", err.Error())
  }

  app.User_db = database.User_Sql_Driver(db)
  app.Git_db = database.Git_Sql_Driver(db)
  app.PAT_db = database.PAT_Sql_Driver(db)

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
