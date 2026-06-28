package main

import (
	"github.com/Ankumeah/SourceFlameBackend/internal/apis"
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/middlewares"

	"github.com/gin-gonic/gin"

	"context"
	"log"
	"sync"
	"time"
)

var Ctx = context.Background()
var app = &a.App{Settings: &a.Settings{}}

func main() {
	load_env(app.Settings)
	app.Settings.PAT_PREFIX = "sf_"
	app.Settings.PAT_LENGTH = 32
	app.Settings.AUTH_CHALLENGE_HEADER = `Basic realm="SourceFlame"`
	app.Settings.JWT_LEEWAY = time.Second * 10
	app.Settings.TOKEN_LENGTH = 64
	app.Settings.TOKEN_NAMESPACE = "refresh:"

	var wg sync.WaitGroup
	wg.Go(func() { connect_session_store(app) })
	wg.Go(func() { connect_database(app) })
	wg.Go(func() { get_handlers(app) })
	wg.Wait()

	log.Println("Starting http server")
	r := gin.Default()
	apiGroup := r.Group(
		"/api/"+app.Settings.API_VERSION+"/",
		middlewars.Log_Middleware(),
	)
	apis.Apis(apiGroup, app)

	log.Println("Running backend on port: " + app.Settings.BACKEND_PORT)
	log.Println("API_VERSION: " + app.Settings.API_VERSION)

	r.Run(":" + app.Settings.BACKEND_PORT)
}
