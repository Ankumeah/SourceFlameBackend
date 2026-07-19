package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"
)

var Ctx = context.Background()
var envVars = map[string]string{
	"DATABASE_URL":            "",
	"DATABASE_MAX_CONNS":      "",
	"DATABASE_MAX_IDLE_CONNS": "",
	"DATABASE_MAX_LIFETIME":   "",
	"DATABASE_MAX_IDLE_TIME":  "",
}

func loadEnv() {
	for env := range envVars {
		Env, ok := os.LookupEnv(env)
		if !ok {
			log.Fatalf("Unset env var: %v\n", env)
		}

		envVars[env] = Env
	}
}

func main() {
	loadEnv()

	maxConn, err := strconv.Atoi(envVars["DATABASE_MAX_CONNS"])
	if err != nil {
		log.Fatalf("Error while parsing DATABASE_MAX_CONNS: %v\n", err.Error())
	}
	maxIdle, err := strconv.Atoi(envVars["DATABASE_MAX_IDLE_CONNS"])
	if err != nil {
		log.Fatalf("Error while parsing DATABASE_MAX_IDLE_CONNS: %v\n", err.Error())
	}
	maxLifetime, err := time.ParseDuration(envVars["DATABASE_MAX_LIFETIME"])
	if err != nil {
		log.Fatalf("Error while parsing DATABASE_MAX_LIFETIME: %v\n", err.Error())
	}
	maxIdleTime, err := time.ParseDuration(envVars["DATABASE_MAX_IDLE_TIME"])
	if err != nil {
		log.Fatalf("Error while parsing DATABASE_MAX_IDLE_TIME: %v\n", err.Error())
	}
	config := NewSqlConfig(int(maxConn), int(maxIdle), maxLifetime, maxIdleTime)

	db, err := GetDBConnection(Ctx, envVars["DATABASE_URL"], config)
	if err != nil {
		log.Fatalf("Error while connecting to database: %v\n", err.Error())
	}
	if err = ExecSchema(Ctx, db); err != nil {
		log.Fatalf("Error while execing schema: %v\n", err.Error())
	}

	log.Println("Connected to database")
}
