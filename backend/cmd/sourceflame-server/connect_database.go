package main

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/database"

	"log"
)

func connectDatabase(app *a.App) {
	log.Println("Connecting to database")

	config := database.NewSqlConfig(
		app.Settings.DATABASE_MAX_CONNS,
		app.Settings.DATABASE_MAX_IDLE_CONNS,
		app.Settings.DATABASE_MAX_LIFETIME,
		app.Settings.DATABASE_MAX_IDLE_TIME,
	)

	db, err := database.GetDBConnection(Ctx, app.Settings.DATABASE_URL, config)
	if err != nil {
		log.Fatalf("Error while connecting to database: %v\n", err.Error())
	}

	app.UserDb = database.UserSqlDriver(db)
	app.GitDb = database.GitSqlDriver(db)
	app.PATDb = database.PATSqlDriver(db)

	log.Println("Connected to database")
}
