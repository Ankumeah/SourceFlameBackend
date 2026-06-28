package main

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/database"

	"log"
)

func connect_database(app *a.App) {
	log.Println("Connecting to database")

	config := database.New_Sql_Config(
		app.Settings.DATABASE_MAX_CONNS,
		app.Settings.DATABASE_MAX_IDLE_CONNS,
		app.Settings.DATABASE_MAX_LIFETIME,
		app.Settings.DATABASE_MAX_IDLE_TIME,
	)

	db, err := database.Get_DB_Connection(Ctx, app.Settings.DATABASE_URL, config)
	if err != nil {
		log.Fatalf("Error while connecting to database: %v\n", err.Error())
	}

	app.User_db = database.User_Sql_Driver(db)
	app.Git_db = database.Git_Sql_Driver(db)
	app.PAT_db = database.PAT_Sql_Driver(db)

	log.Println("Connected to database")
}
