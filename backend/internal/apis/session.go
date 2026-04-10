package apis

import (
	a "github.com/Ankumeah/DeltaBase/internal/app"
	"github.com/Ankumeah/DeltaBase/internal/jwt"
	"github.com/Ankumeah/DeltaBase/internal/middlewares"
	"github.com/Ankumeah/DeltaBase/internal/session_store"

	"github.com/gin-gonic/gin"

	"net/http"
)

func session(r *gin.RouterGroup, app *a.App) {
  group := r.Group("/session", middlewars.Verify_Session_Middleware(app.Store))

  group.POST("/renew_jwt", func (c *gin.Context) {
    username := c.GetString("username")

    token, err := jwt.Issue_jwt(username)
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
    username := c.GetHeader("username")
    session := c.GetHeader("session")

    token, err := app.Store.Add_Session(ctx, username)
    if err == session_store.Error_too_many_tokens {
      c.JSON(http.StatusForbidden, gin.H { "error": "Too many tokens" })
      return
    } else if err != nil {
      c.JSON(internal_server_error())
      return
    } else {
      c.JSON(http.StatusOK, gin.H { "refresh_token": token })
    }

    if err = app.Store.Delete_Session(ctx, username, session); err != nil {
      c.JSON(internal_server_error())
    }
  })

  group.DELETE("", func (c *gin.Context) {
    ctx := c.Request.Context()
    session := c.GetHeader("session")
    username := c.GetHeader("username")

    if err := app.Store.Delete_Session(ctx, username, session); err != nil {
      c.JSON(internal_server_error())
      return
    }

    c.Status(http.StatusOK)
  })
}
