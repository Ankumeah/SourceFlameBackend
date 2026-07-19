package main

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"

	"log"
	"os"
	"strconv"
	"time"
)

var envVars = map[string]string{
	"API_VERSION":  "",
	"BACKEND_PORT": "",

	"SESSION_STORE_SESSIONS_USERNAME": "",
	"SESSION_STORE_SESSIONS_PASSWORD": "",
	"SESSION_STORE_HOSTNAME":          "",
	"SESSION_STORE_PORT":              "",
	"SESSION_STORE_TYPE":              "",

	"DATABASE_URL":            "",
	"DATABASE_MAX_CONNS":      "",
	"DATABASE_MAX_IDLE_CONNS": "",
	"DATABASE_MAX_LIFETIME":   "",
	"DATABASE_MAX_IDLE_TIME":  "",

	"JWT_KEY":      "",
	"JWT_LIFESPAN": "",

	"TOKEN_TIMEOUT": "",
	"TOKEN_LIMIT":   "",
}

func loadEnv(s *a.Settings) {
	log.Println("Loading env")

	for env := range envVars {
		_env, ok := os.LookupEnv(env)
		if !ok {
			log.Fatalf("Unset env var: %v\n", env)
		}

		envVars[env] = _env
	}

	setSettings(s)

	log.Println("Loaded env")
}

func setSettings(s *a.Settings) {
	s.API_VERSION = envVars["API_VERSION"]
	s.BACKEND_PORT = envVars["BACKEND_PORT"]

	s.SESSION_STORE_SESSIONS_USERNAME = envVars["SESSION_STORE_SESSIONS_USERNAME"]
	s.SESSION_STORE_SESSIONS_PASSWORD = envVars["SESSION_STORE_SESSIONS_PASSWORD"]
	s.SESSION_STORE_HOSTNAME = envVars["SESSION_STORE_HOSTNAME"]
	s.SESSION_STORE_PORT = envVars["SESSION_STORE_PORT"]
	s.SESSION_STORE_TYPE = envVars["SESSION_STORE_TYPE"]

	s.DATABASE_URL = envVars["DATABASE_URL"]

	var err error
	s.DATABASE_MAX_CONNS, err = strconv.Atoi(envVars["DATABASE_MAX_CONNS"])
	if err != nil {
		log.Fatalf("Error while parsing DATABASE_MAX_CONNS: %v\n", err.Error())
	}
	s.DATABASE_MAX_IDLE_CONNS, err = strconv.Atoi(envVars["DATABASE_MAX_IDLE_CONNS"])
	if err != nil {
		log.Fatalf("Error while parsing DATABASE_MAX_IDLE_CONNS: %v\n", err.Error())
	}
	s.DATABASE_MAX_LIFETIME, err = time.ParseDuration(envVars["DATABASE_MAX_LIFETIME"])
	if err != nil {
		log.Fatalf("Error while parsing DATABASE_MAX_LIFETIME: %v\n", err.Error())
	}
	s.DATABASE_MAX_IDLE_TIME, err = time.ParseDuration(envVars["DATABASE_MAX_IDLE_TIME"])
	if err != nil {
		log.Fatalf("Error while parsing DATABASE_MAX_IDLE_TIME: %v\n", err.Error())
	}

	s.JWT_KEY = []byte(envVars["JWT_KEY"])
	s.JWT_LIFESPAN, err = time.ParseDuration(envVars["JWT_LIFESPAN"])
	if err != nil {
		log.Fatalf("Error while parsing JWT_LIFESPAN: %v\n", err.Error())
	}

	s.TOKEN_TIMEOUT, err = time.ParseDuration(envVars["TOKEN_TIMEOUT"])
	if err != nil {
		log.Fatalf("Error while parsing TOKEN_TIMEOUT: %v\n", err.Error())
	}
	s.TOKEN_LIMIT, err = strconv.Atoi(envVars["TOKEN_LIMIT"])
	if err != nil {
		log.Fatalf("Error while parsing TOKEN_LIMIT: %v\n", err.Error())
	}
}
