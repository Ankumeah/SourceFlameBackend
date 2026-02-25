package apis

import (
	"github.com/Ankumeah/DeltaBase/internal/database"
	"github.com/Ankumeah/DeltaBase/internal/jwt"
	"github.com/Ankumeah/DeltaBase/internal/session_store"

	"github.com/gin-gonic/gin"

	"net/http"
)

func login(r *gin.RouterGroup, store *session_store.Session_store, db *database.Database) {
  r.POST("/login", func(c *gin.Context) {
    username, ok := get_id(c, db)
    if !ok { return }

    add_session(c, store, username)
  })
}

func get_id(c *gin.Context, db *database.Database) (string, bool) {
  ctx := c.Request.Context()
  var request struct {
    username string
    password string
  }

  if err := c.ShouldBindJSON(&request); err != nil {
    c.JSON(http.StatusBadRequest, gin.H { "error": err.Error() })
    return "", false
  }

  _, err := db.Get_Id(ctx, request.username)
  if err == database.Error_invalid_user {
    _, err = db.Add_User(ctx, request.username, request.password)
  }
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server error" })
    return "", false
  }

  return request.username, true
}

func add_session(c *gin.Context, store *session_store.Session_store, username string) bool {
  ctx := c.Request.Context()

  token, err := store.Add_Session(ctx, username)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H { "error": "Invalid server error" })
    return false
  }

  new_jwt, err := jwt.Issue_jwt(username)
  if err != nil {
    store.Delete_Session(ctx, token)

    c.JSON(http.StatusInternalServerError, gin.H { "error": "Invalid server error" })
    return false
  } else {
    c.JSON(http.StatusOK, gin.H { "JWT": new_jwt, "refresh_token": token })
    return  true
  }
}
