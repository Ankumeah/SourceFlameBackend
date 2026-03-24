package apis

import (
	"github.com/Ankumeah/DeltaBase/internal/jwt"
	"github.com/Ankumeah/DeltaBase/internal/session_store"
	"github.com/Ankumeah/DeltaBase/internal/middlewears"

	"github.com/gin-gonic/gin"

	"net/http"
  "fmt"
)

func session(r *gin.RouterGroup, store *session_store.Session_store) {
  group := r.Group("/session", middlewears.Verify_Session_Middlewear(store))

  group.POST("/renew_jwt", func (c *gin.Context) {
    username, _ := c.Get("username")

    token, err := jwt.Issue_jwt(fmt.Sprint(username))
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server error" })
      return
    } else {
      c.JSON(http.StatusOK, gin.H { "JWT": token })
      return
    }
  })

  group.POST("/renew_session", func (c *gin.Context) {
    ctx := c.Request.Context()
    username := c.GetHeader("username")
    session := c.GetHeader("session")

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

    if err = store.Delete_Session(ctx, username, session); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server error" })
    }
  })

  group.POST("/logout", func (c *gin.Context) {
    ctx := c.Request.Context()
    session := c.GetHeader("session")
    username := c.GetHeader("username")

    if err := store.Delete_Session(ctx, username, session); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server error" })
      return
    }

    c.Status(http.StatusOK)
  })
}
