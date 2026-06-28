package main

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"

	"log"
	"os"
	"strconv"
	"time"
)

var env_vars = map[string]string{
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

func load_env(s *a.Settings) {
	log.Println("Loading env")

	for env := range env_vars {
		_env, ok := os.LookupEnv(env)
		if !ok {
			log.Fatalf("Unset env var: %v\n", env)
		}

		env_vars[env] = _env
	}

	set_settings(s)

	log.Println("Loaded env")
}

func set_settings(s *a.Settings) {
	s.API_VERSION = env_vars["API_VERSION"]
	s.BACKEND_PORT = env_vars["BACKEND_PORT"]

	s.SESSION_STORE_SESSIONS_USERNAME = env_vars["SESSION_STORE_SESSIONS_USERNAME"]
	s.SESSION_STORE_SESSIONS_PASSWORD = env_vars["SESSION_STORE_SESSIONS_PASSWORD"]
	s.SESSION_STORE_HOSTNAME = env_vars["SESSION_STORE_HOSTNAME"]
	s.SESSION_STORE_PORT = env_vars["SESSION_STORE_PORT"]
	s.SESSION_STORE_TYPE = env_vars["SESSION_STORE_TYPE"]

	s.DATABASE_URL = env_vars["DATABASE_URL"]

	var err error
	s.DATABASE_MAX_CONNS, err = strconv.Atoi(env_vars["DATABASE_MAX_CONNS"])
	if err != nil {
		log.Fatalf("Error while parseing DATABASE_MAX_CONNS: %v\n", err.Error())
	}
	s.DATABASE_MAX_IDLE_CONNS, err = strconv.Atoi(env_vars["DATABASE_MAX_IDLE_CONNS"])
	if err != nil {
		log.Fatalf("Error while parseing DATABASE_MAX_IDLE_CONNS: %v\n", err.Error())
	}
	s.DATABASE_MAX_LIFETIME, err = time.ParseDuration(env_vars["DATABASE_MAX_LIFETIME"])
	if err != nil {
		log.Fatalf("Error while parseing DATABASE_MAX_LIFETIME: %v\n", err.Error())
	}
	s.DATABASE_MAX_IDLE_TIME, err = time.ParseDuration(env_vars["DATABASE_MAX_IDLE_TIME"])
	if err != nil {
		log.Fatalf("Error while parseing DATABASE_MAX_IDLE_TIME: %v\n", err.Error())
	}

	s.JWT_KEY = []byte(env_vars["JWT_KEY"])
	s.JWT_LIFESPAN, err = time.ParseDuration(env_vars["JWT_LIFESPAN"])
	if err != nil {
		log.Fatalf("Error while parseing JWT_LIFESPAN: %v\n", err.Error())
	}

	s.TOKEN_TIMEOUT, err = time.ParseDuration(env_vars["TOKEN_TIMEOUT"])
	if err != nil {
		log.Fatalf("Error while parseing TOKEN_TIMEOUT: %v\n", err.Error())
	}
	s.TOKEN_LIMIT, err = strconv.Atoi(env_vars["TOKEN_LIMIT"])
	if err != nil {
		log.Fatalf("Error while parseing TOKEN_LIMIT: %v\n", err.Error())
	}
}
