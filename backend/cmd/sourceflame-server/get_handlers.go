package main

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/jwt"
	"github.com/Ankumeah/SourceFlameBackend/internal/pat"
)

func getHandlers(app *a.App) {
	app.PATHandler = pat.GetPATHandler(
		app.Settings.PAT_PREFIX,
		app.Settings.PAT_LENGTH,
	)
	app.JWTHandler = jwt.GetJWTHandler(
		app.Settings.JWT_LEEWAY,
		app.Settings.JWT_LIFESPAN,
		app.Settings.JWT_KEY,
	)
}
