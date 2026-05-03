package middlewars

import (
	"github.com/Ankumeah/DeltaBase/internal/database"

	"github.com/gin-gonic/gin"

	"net/http"
)

func Self_Login_Middleware(user_db *database.User_db) gin.HandlerFunc {
  return func(c *gin.Context) {
    ctx := c.Request.Context()
    var request struct {
      Username string `json:"username" binding:"required"`
      Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
      c.JSON(http.StatusBadRequest, gin.H { "error": err.Error() })
      return
    }

    user_id, err := user_db.Get_Id(ctx, request.Username)
    if err == database.Error_invalid_user {
      c.JSON(http.StatusBadRequest, gin.H { "error": "Invalid user" })
      c.Abort()
      return
    } else if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server error" })
      c.Abort()
      return
    }

    valid, err := user_db.Verify_User(ctx, user_id, request.Password)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server error" })
      c.Abort()
      return
    } else if !valid {
      c.JSON(http.StatusUnauthorized, gin.H { "error": "Invalid password" })
      c.Abort()
      return
    } else {
      c.Set("user_id", user_id)
      c.Next()
      return
    }
  }
}
