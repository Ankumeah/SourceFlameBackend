package app

import (
	"github.com/Ankumeah/SourceFlameBackend/internal/database"
	"github.com/Ankumeah/SourceFlameBackend/internal/pat"
	"github.com/Ankumeah/SourceFlameBackend/internal/session_store"
)

type App struct {
  User_db *database.User_db
  Git_db *database.Git_db
  PAT_db *database.PAT_db
  Store *session_store.Session_store
  PAT_Handler *pat.PAT_Handler
}
