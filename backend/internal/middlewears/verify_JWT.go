package middlewears

import (
	"github.com/Ankumeah/DeltaBase/internal/jwt"

	"github.com/gin-gonic/gin"

  "net/http"
)

func Verify_JWT_Middlewear() gin.HandlerFunc {
  return func(c *gin.Context) {
    token := c.GetHeader("JWT")
    username := c.GetHeader("username")

    valid, err := jwt.Validate_jwt(username, token)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server error" })
      c.Abort()
      return
    } else if !valid {
      c.JSON(http.StatusUnauthorized, gin.H { "error": "Invalid JWT" })
      c.Abort()
      return
    } else {
      c.Set("username", username)
      c.Next()
      return
    }
  }
}
