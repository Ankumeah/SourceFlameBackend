package apis

import (
	"github.com/Ankumeah/DeltaBase/internal/session_store"

	"github.com/gin-gonic/gin"

	"net/http"
)

func logout(r *gin.RouterGroup, store *session_store.Session_store) {
  r.POST("/logout", func (c *gin.Context) {
    ctx := c.Request.Context()
    session := c.GetHeader("session")

    if err := store.Delete_Session(ctx, session); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server error" })
    }
  })
}
