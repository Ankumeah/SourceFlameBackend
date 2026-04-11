package main

import (
	"github.com/Ankumeah/DeltaBase/internal/apis"
	a "github.com/Ankumeah/DeltaBase/internal/app"
	"github.com/Ankumeah/DeltaBase/internal/database"
	"github.com/Ankumeah/DeltaBase/internal/middlewares"
	"github.com/Ankumeah/DeltaBase/internal/session_store"

	"github.com/gin-gonic/gin"

	"context"
	"log"
	"os"
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

var app *a.App

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
  conn_config := session_store.Connection_Config {
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
  log.Println("Connected to session store")
}

func connect_database() {
  conn_config := database.Connection_Config {
    Username: env_vars["DATABASE_SESSIONS_USERNAME"],
    Password: env_vars["DATABASE_SESSIONS_PASSWORD"],
    Hostname: env_vars["DATABASE_HOSTNAME"],
    Port: env_vars["DATABASE_PORT"],
    Db_name: env_vars["DATABASE_DB"],
    Db_config: env_vars["DATABASE_CONFIG"],
  }

  user_db, err := database.User_Postgres_Driver(Ctx, conn_config)
  if err != nil {
    log.Fatalf("Error while connecting to database: %v\n", err.Error())
  }

  git_db, err := database.Git_Postgres_Driver(Ctx, conn_config)
  if err != nil {
    log.Fatalf("Error while connecting to database: %v\n", err.Error())
  }

  app.User_db = user_db
  app.Git_db = git_db
  log.Println("Connected to database")
}

func main() {
	log.Println("Loading env")
  load_env()

	log.Println("Connecting to DBs")
  connect_session_store()
  connect_database()

	log.Println("String http server")
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
