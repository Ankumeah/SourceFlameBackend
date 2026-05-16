package middlewars

import (
	"github.com/Ankumeah/SourceFlameBackend/internal/session_store"

	"github.com/gin-gonic/gin"

  "net/http"
  "log"
)

func Verify_Session_Middleware(store *session_store.Session_store) gin.HandlerFunc {
  return func(c *gin.Context) {
    ctx := c.Request.Context()

    session := c.GetHeader("X-Authorization")
    username := c.GetHeader("X-Username")

    valid, err := store.Validate_Session(ctx, username, session)
    if err != nil {
      log.Printf("Error while validateing session: %v\n", err.Error())
      c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server error" })
      c.Abort()
      return
    } else if !valid {
      c.JSON(http.StatusUnauthorized, gin.H { "error": "Invalid refresh token" })
      c.Abort()
      return
    } else {
      c.Set("username", username)
      c.Set("session", session)
      c.Next()
      return
    }
  }
}
