package apis

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/middlewares"
	"github.com/Ankumeah/SourceFlameBackend/internal/session_store"

	"github.com/gin-gonic/gin"

	"net/http"
)

func session(r *gin.RouterGroup, app *a.App) {
  group := r.Group("/session", middlewars.Verify_Session_Middleware(app.Store))

  group.POST("/renew_jwt", func (c *gin.Context) {
    username := c.GetString("username")

    token, err := app.JWT_Handler.Issue_jwt(username)
    if err != nil {
      c.JSON(internal_server_error())
      return
    } else {
      c.JSON(http.StatusOK, gin.H { "JWT": token })
      return
    }
  })

  group.POST("/renew_session", func (c *gin.Context) {
    ctx := c.Request.Context()
    username := c.GetString("username")
    session := c.GetString("session")

    token, err := app.Store.Add_Session(ctx, username)
    if err == session_store.Error_too_many_tokens {
      c.JSON(http.StatusForbidden, gin.H { "error": "Too many tokens" })
      return
    } else if err != nil {
      c.JSON(internal_server_error())
      return
    }

    if err = app.Store.Delete_Session(ctx, username, session); err != nil {
      c.JSON(internal_server_error())
      app.Store.Delete_Session(ctx, username, token)
      return
    }

    c.JSON(http.StatusOK, gin.H { "refresh_token": token })
  })

  group.DELETE("", func (c *gin.Context) {
    ctx := c.Request.Context()
    session := c.GetString("session")
    username := c.GetString("username")

    if err := app.Store.Delete_Session(ctx, username, session); err != nil {
      c.JSON(internal_server_error())
      return
    }

    c.Status(http.StatusOK)
  })
}
