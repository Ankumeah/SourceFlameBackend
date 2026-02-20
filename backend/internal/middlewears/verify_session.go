package middlewears

import (
	"github.com/Ankumeah/DeltaBase/internal/session_store"

	"github.com/gin-gonic/gin"

  "net/http"
)

func Verify_Session_Middlewear(store *session_store.Session_store) gin.HandlerFunc {
  return func(c *gin.Context) {
    ctx := c.Request.Context()

    session := c.GetHeader("session")
    username := c.GetHeader("username")

    valid, err := store.Validate_Session(ctx, username, session)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server error" })
      c.Abort()
      return
    } else if !valid {
      c.JSON(http.StatusUnauthorized, gin.H { "error": "Invalid refresh token" })
      c.Abort()
      return
    } else {
      c.Next()
      return
    }
  }
}
