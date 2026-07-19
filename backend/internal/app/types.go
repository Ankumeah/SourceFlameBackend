package app

import (
	"github.com/Ankumeah/SourceFlameBackend/internal/database"
	"github.com/Ankumeah/SourceFlameBackend/internal/jwt"
	"github.com/Ankumeah/SourceFlameBackend/internal/pat"
	"github.com/Ankumeah/SourceFlameBackend/internal/session_store"

	"time"
)

type App struct {
	UserDb     *database.UserDb
	GitDb      *database.GitDb
	PATDb      *database.PATDb
	Store      *session_store.SessionStore
	PATHandler *pat.PATHandler
	JWTHandler *jwt.JWTHandler
	Settings   *Settings
}

type Settings struct {
	API_VERSION  string
	BACKEND_PORT string

	SESSION_STORE_SESSIONS_USERNAME string
	SESSION_STORE_SESSIONS_PASSWORD string
	SESSION_STORE_HOSTNAME          string
	SESSION_STORE_PORT              string
	SESSION_STORE_TYPE              string

	DATABASE_URL            string
	DATABASE_MAX_CONNS      int
	DATABASE_MAX_IDLE_CONNS int
	DATABASE_MAX_LIFETIME   time.Duration
	DATABASE_MAX_IDLE_TIME  time.Duration

	JWT_KEY      []byte
	JWT_LIFESPAN time.Duration
	JWT_LEEWAY   time.Duration

	TOKEN_TIMEOUT   time.Duration
	TOKEN_LIMIT     int
	TOKEN_LENGTH    int
	TOKEN_NAMESPACE string

	PAT_PREFIX string
	PAT_LENGTH uint8

	AUTH_CHALLENGE_HEADER string
}
