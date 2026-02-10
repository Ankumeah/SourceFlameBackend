package main

import (
	"github.com/Ankumeah/DeltaBase/internal/apis"
  "github.com/Ankumeah/DeltaBase/internal/session_store"

	"github.com/gin-gonic/gin"

	"fmt"
  "os"
  "context"
)

var env_vars = map[string]string {
  "BACKEND_PORT": "",
  "SESSION_STORE_USERNAME": "",
  "SESSION_STORE_PASSWORD": "",
  "SESSION_STORE_HOSTNAME": "",
  "SESSION_STORE_PORT": "",
  "SESSION_STORE_TYPE": "",
  "JWT_KEY": "",
}
var Ctx = context.Background()

var store *session_store.Session_store

func init() {
  for env := range env_vars {
    _env, ok := os.LookupEnv(env)
    if !ok {
      panic("env var " + env + " not set")
    }

    env_vars[env] = _env
  }

  session_store_url := fmt.Sprintf("%v://%v:%v@%v:%v",
    env_vars["SESSION_STORE_TYPE"],
    env_vars["SESSION_STORE_USERNAME"],
    env_vars["SESSION_STORE_PASSWORD"],
    env_vars["SESSION_STORE_HOSTNAME"],
    env_vars["SESSION_STORE_PORT"],
  )

  _store, err := session_store.Get_Redis_Driver(Ctx, session_store_url)
  if err != nil {
    panic("Error while connecting to session store: " + err.Error())
  }

  store = _store
}

func main() {
	fmt.Println("Starting http server")

	r := gin.Default()

	apiGroup := r.Group("/api")
	apis.Apis(apiGroup, store)

  port, ok := os.LookupEnv("BACKEND_PORT")
  if !ok {
    panic("env var BACKEND_PORT not set")
  }

  fmt.Println("Running backend on port: " + port)
	r.Run(":" + port)
}
