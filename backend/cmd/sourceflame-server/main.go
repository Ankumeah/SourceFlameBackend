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
	loadEnv(app.Settings)
	app.Settings.PAT_PREFIX = "sf_"
	app.Settings.PAT_LENGTH = 32
	app.Settings.AUTH_CHALLENGE_HEADER = `Basic realm="SourceFlame"`
	app.Settings.JWT_LEEWAY = time.Second * 10
	app.Settings.TOKEN_LENGTH = 64
	app.Settings.TOKEN_NAMESPACE = "refresh:"

	var wg sync.WaitGroup
	wg.Go(func() { connectSessionStore(app) })
	wg.Go(func() { connectDatabase(app) })
	wg.Go(func() { getHandlers(app) })
	wg.Wait()

	log.Println("Starting http server")
	r := gin.Default()
	apiGroup := r.Group(
		"/api/"+app.Settings.API_VERSION+"/",
		middlewares.LogMiddleware(),
	)
	apis.Apis(apiGroup, app)

	log.Println("Running backend on port: " + app.Settings.BACKEND_PORT)
	log.Println("API_VERSION: " + app.Settings.API_VERSION)

	err := r.Run(":" + app.Settings.BACKEND_PORT)
	if err != nil {
		log.Fatalf("Error while running server: %v\n", err.Error())
	}
}
