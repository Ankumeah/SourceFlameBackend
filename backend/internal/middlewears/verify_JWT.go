package middlewears

import (
	"github.com/Ankumeah/DeltaBase/internal/jwt"

	"github.com/gin-gonic/gin"

  "net/http"
)

func Verify_JWT_Middlewear() gin.HandlerFunc {
  return func(c *gin.Context) {
    token := c.GetHeader("JWT")

    username, err := jwt.Validate_jwt(token)
    if err == jwt.Error_invalid_JWT {
      c.JSON(http.StatusUnauthorized, gin.H { "error": "Invalid JWT" })
      c.Abort()
      return
    } else if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server error" })
      c.Abort()
      return
    } else {
      c.Set("username", username)
      c.Next()
      return
    }
  }
}
