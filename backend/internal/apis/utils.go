package apis

import (
  "github.com/Ankumeah/SourceFlameBackend/internal/database"

	"github.com/gin-gonic/gin"

  "net/http"
)

func get_user_id(c *gin.Context, user_db *database.User_db, username string) (uint64, bool) {
  ctx := c.Request.Context()

  user_id, err := user_db.Get_Id(ctx, username)
  if err == database.Error_invalid_user {
    c.JSON(http.StatusBadRequest, gin.H { "error": "Invalid user" })
    return 0, false
  } else if err != nil {
    c.JSON(internal_server_error())
    return 0, false
  } else {
    return user_id, true
  }
}
