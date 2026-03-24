package apis

import (
	"github.com/Ankumeah/DeltaBase/internal/database"

	"github.com/gin-gonic/gin"

	"net/http"
)

func user(r *gin.RouterGroup, user_db *database.User_db) {
  group := r.Group("/user")

  group.GET("/:username/*action", func(c *gin.Context) {
    switch c.Param("action") {
      case "/creation":
        creation(c, user_db)
    }
  })
}

func creation(c *gin.Context, user_db *database.User_db) {
  ctx := c.Request.Context()
  username := c.Query("username")

  user_id, err := user_db.Get_Id(ctx, username)
  if err == database.Error_invalid_user {
    c.JSON(http.StatusBadRequest, gin.H { "error": "Unknown user" })
    return
  } else if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server error" })
    return
  }

  timestamp, err := user_db.Get_Creation(ctx, user_id)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server error" })
    return
  } else {
    c.JSON(http.StatusOK, gin.H { "timestamp": timestamp })
  }
}
