package main

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/pat"
	"github.com/Ankumeah/SourceFlameBackend/internal/jwt"
)

func get_handlers(app *a.App) {
  app.PAT_Handler = pat.Get_PAT_Handler(
    app.Settings.PAT_PREFIX,
    app.Settings.PAT_LENGTH,
  )
  app.JWT_Handler = jwt.Get_JWT_Handler(
    app.Settings.JWT_LEEWAY,
    app.Settings.JWT_LIFESPAN,
    app.Settings.JWT_KEY,
  )
}
