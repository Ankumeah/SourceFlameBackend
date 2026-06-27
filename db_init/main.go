package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"
)

var Ctx = context.Background()
var env_vars = map[string]string{
	"DATABASE_URL":            "",
	"DATABASE_MAX_CONNS":      "",
	"DATABASE_MAX_IDLE_CONNS": "",
	"DATABASE_MAX_LIFETIME":   "",
	"DATABASE_MAX_IDLE_TIME":  "",
}

func load_env() {
	for env := range env_vars {
		_env, ok := os.LookupEnv(env)
		if !ok {
			log.Fatalf("Unset env var: %v\n", env)
		}

		env_vars[env] = _env
	}
}

func main() {
	load_env()

	max_conn, err := strconv.Atoi(env_vars["DATABASE_MAX_CONNS"])
	if err != nil {
		log.Fatalf("Error while parsing DATABASE_MAX_CONNS: %v\n", err.Error())
	}
	max_idle, err := strconv.Atoi(env_vars["DATABASE_MAX_IDLE_CONNS"])
	if err != nil {
		log.Fatalf("Error while parsing DATABASE_MAX_IDLE_CONNS: %v\n", err.Error())
	}
	max_lifetime, err := time.ParseDuration(env_vars["DATABASE_MAX_LIFETIME"])
	if err != nil {
		log.Fatalf("Error while parsing DATABASE_MAX_LIFETIME: %v\n", err.Error())
	}
	max_idle_time, err := time.ParseDuration(env_vars["DATABASE_MAX_IDLE_TIME"])
	if err != nil {
		log.Fatalf("Error while parsing DATABASE_MAX_IDLE_TIME: %v\n", err.Error())
	}
	config := New_Sql_Config(int(max_conn), int(max_idle), max_lifetime, max_idle_time)

	db, err := Get_DB_Connection(Ctx, env_vars["DATABASE_URL"], config)
	if err != nil {
		log.Fatalf("Error while connecting to database: %v\n", err.Error())
	}
	if err = Exec_Schema(Ctx, db); err != nil {
		log.Fatalf("Error while execing schema: %v\n", err.Error())
	}

	log.Println("Connected to database")
}
