package app

import (
  "github.com/Ankumeah/DeltaBase/internal/database"
  "github.com/Ankumeah/DeltaBase/internal/session_store"
)

type App struct {
  User_db *database.User_db
  Git_db *database.Git_db
  Store *session_store.Session_store
}
