package apis

import (
	"github.com/Ankumeah/DeltaBase/internal/session_store"
	"github.com/Ankumeah/DeltaBase/internal/database"

	"github.com/gin-gonic/gin"
)

func Apis(
  r *gin.RouterGroup,
  store *session_store.Session_store,
  user_db *database.User_db,
  git_db *database.Git_db,
) {
  login(r, store, user_db)
  repos(r, git_db)
  session(r, store)
  user(r, user_db)
}
