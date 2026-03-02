package main

import (
	"github.com/Ankumeah/DeltaBase/internal/apis"
  "github.com/Ankumeah/DeltaBase/internal/session_store"
  "github.com/Ankumeah/DeltaBase/internal/middlewears"
  "github.com/Ankumeah/DeltaBase/internal/database"

	"github.com/gin-gonic/gin"

	"log"
  "os"
  "context"
)

var env_vars = map[string]string {
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
  "JWT_KEY": "",
}
var Ctx = context.Background()

var store *session_store.Session_store
var db *database.Database

func init() {
  load_env()
}

func load_env() {
  for env := range env_vars {
    _env, ok := os.LookupEnv(env)
    if !ok {
      panic("env var " + env + " not set")
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
    panic("Error while connecting to session store: " + err.Error())
  }

  store = _store
  log.Println("Connected to session store")
}

func connect_database() {
  _db, err := database.Get_Postgres_Driver(Ctx,
    env_vars["DATABASE_USER"],
    env_vars["DATABASE_PASSWORD"],
    env_vars["DATABASE_HOST"],
    env_vars["DATABASE_PORT"],
    env_vars["DATABASE_DB"],
  )
  if err != nil {
    panic("Error while connecting to database: " + err.Error())
  }

  db = _db
  log.Println("Connected to database")
}

func main() {
	log.Println("Starting http server")

  connect_session_store()
  connect_database()

	r := gin.Default()

	apiGroup := r.Group("/api",
    middlewears.Log_Middlewear(),
  )
	apis.Apis(apiGroup, store ,db)

  jwtRequiredGroup := apiGroup.Group("/jwt",
    middlewears.Verify_JWT_Middlewear(),
  )
  apis.JWT_Needed_Apis(jwtRequiredGroup)

  sessionRequiredGroup := apiGroup.Group("/session",
    middlewears.Verify_Session_Middlewear(store),
  )
  apis.Session_Needed_Apis(sessionRequiredGroup, store)

  log.Println("Running backend on port: " + env_vars["BACKEND_PORT"])
	r.Run(":" + env_vars["BACKEND_PORT"])
}
