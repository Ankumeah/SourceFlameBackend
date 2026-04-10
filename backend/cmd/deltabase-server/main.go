package main

import (
	"github.com/Ankumeah/DeltaBase/internal/apis"
  "github.com/Ankumeah/DeltaBase/internal/session_store"
  "github.com/Ankumeah/DeltaBase/internal/middlewares"
  "github.com/Ankumeah/DeltaBase/internal/database"

	"github.com/gin-gonic/gin"

	"log"
  "os"
  "context"
)

var env_vars = map[string]string {
  "API_VERSION": "",
  "BACKEND_PORT": "",
  "SESSION_STORE_SESSIONS_USERNAME": "",
  "SESSION_STORE_SESSIONS_PASSWORD": "",
  "SESSION_STORE_HOSTNAME": "",
  "SESSION_STORE_PORT": "",
  "DATABASE_USER": "",
  "DATABASE_PASSWORD": "",
  "DATABASE_HOST": "",
  "DATABASE_PORT": "",
  "DATABASE_DB": "",
  "DATABASE_CONFIG": "",
}
var Ctx = context.Background()

var store *session_store.Session_store
var user_db *database.User_db
var git_db *database.Git_db

func init() {
  for env := range env_vars {
    _env, ok := os.LookupEnv(env)
    if !ok {
      log.Fatalf("Unset env var: %v\n", env)
    }

    env_vars[env] = _env
  }
}

func connect_session_store() {
  _store, err := session_store.Get_Redis_Cluster_Driver(Ctx,
    env_vars["SESSION_STORE_SESSIONS_USERNAME"],
    env_vars["SESSION_STORE_SESSIONS_PASSWORD"],
    env_vars["SESSION_STORE_HOSTNAME"],
    env_vars["SESSION_STORE_PORT"],
  )
  if err != nil {
    log.Fatalf("Error while connecting to session store: %v\n", err.Error())
  }

  store = _store
  log.Println("Connected to session store")
}

func connect_database() {
  _user_db, err := database.User_Postgres_Driver(Ctx,
    env_vars["DATABASE_USER"],
    env_vars["DATABASE_PASSWORD"],
    env_vars["DATABASE_HOST"],
    env_vars["DATABASE_PORT"],
    env_vars["DATABASE_DB"],
    env_vars["DATABASE_CONFIG"],
  )
  if err != nil {
    log.Fatalf("Error while connecting to database: %v\n", err.Error())
  }

  _git_db, err := database.Git_Postgres_Driver(Ctx,
    env_vars["DATABASE_USER"],
    env_vars["DATABASE_PASSWORD"],
    env_vars["DATABASE_HOST"],
    env_vars["DATABASE_PORT"],
    env_vars["DATABASE_DB"],
    env_vars["DATABASE_CONFIG"],
  )
  if err != nil {
    log.Fatalf("Error while connecting to database: %v\n", err.Error())
  }

  user_db = _user_db
  git_db = _git_db
  log.Println("Connected to database")
}

func main() {
	log.Println("Starting http server")

  connect_session_store()
  connect_database()

	r := gin.Default()

	apiGroup := r.Group("/api/" + env_vars["API_VERSION"] + "/",
    middlewars.Log_Middleware(),
  )
	apis.Apis(apiGroup, store ,user_db, git_db)

  log.Println("Running backend on port: " + env_vars["BACKEND_PORT"])
  log.Println("API_VERSION: " + env_vars["API_VERSION"])
	r.Run(":" + env_vars["BACKEND_PORT"])
}
