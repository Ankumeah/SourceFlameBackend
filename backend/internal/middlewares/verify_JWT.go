package middlewars

import (
	"github.com/Ankumeah/SourceFlameBackend/internal/jwt"

	"github.com/gin-gonic/gin"

	"net/http"
)

func Verify_JWT_Middleware(handler *jwt.JWT_Handler) gin.HandlerFunc {
  return func(c *gin.Context) {
    token := c.GetHeader("X-Authorization")

    username, err := handler.Validate_jwt(token)
    if err == jwt.Error_invalid_JWT {
      c.JSON(http.StatusUnauthorized, gin.H { "error": "Invalid JWT" })
      c.Abort()
      return
    } else if err != nil {
      c.JSON(http.StatusBadRequest, gin.H { "error": err.Error() })
      c.Abort()
      return
    } else {
      c.Set("username", username)
      c.Next()
      return
    }
  }
}
