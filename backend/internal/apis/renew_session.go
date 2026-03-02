package apis

import (
	"github.com/Ankumeah/DeltaBase/internal/session_store"

	"github.com/gin-gonic/gin"

	"net/http"
)

func renew_session(r *gin.RouterGroup, store *session_store.Session_store) {
  r.POST("/renew_session", func (c *gin.Context) {
    ctx := c.Request.Context()
    username := c.GetHeader("username")

    token, err := store.Add_Session(ctx, username)
    if err == session_store.Error_too_many_tokens {
      c.JSON(http.StatusForbidden, gin.H { "error": "Too many tokens" })
      return
    } else if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server error" })
      return
    } else {
      c.JSON(http.StatusOK, gin.H { "refresh_token": token })
    }
  })
}
