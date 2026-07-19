package main

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/session_store"

	"log"
)

func connectSessionStore(app *a.App) {
	log.Println("Connecting to session store")

	sessionsConfig := session_store.UniversalRedisConfig{
		Username: app.Settings.SESSION_STORE_SESSIONS_USERNAME,
		Password: app.Settings.SESSION_STORE_SESSIONS_PASSWORD,
		Hostname: app.Settings.SESSION_STORE_HOSTNAME,
		Port:     app.Settings.SESSION_STORE_PORT,
	}
	var client session_store.RedisClient
	var err error

	switch app.Settings.SESSION_STORE_TYPE {
	case "redis_standalone":
		client, err = session_store.GetRedisClient(Ctx, sessionsConfig)
	case "redis_cluster":
		client, err = session_store.GetRedisClusterClient(Ctx, sessionsConfig)
	default:
		log.Fatalf("Unsupported session store type: %v\n", app.Settings.SESSION_STORE_TYPE)
	}
	if err != nil {
		log.Panicf("Error while connecting to session store: %v\n", err.Error())
	}

	app.Store = session_store.GetSessionsUniversalRedisDriver(
		client,
		app.Settings.TOKEN_TIMEOUT,
		app.Settings.TOKEN_LIMIT,
		app.Settings.TOKEN_LENGTH,
		app.Settings.TOKEN_NAMESPACE,
	)
	log.Println("Connected to session store")
}
