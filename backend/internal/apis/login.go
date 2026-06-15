package apis

import (

	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/database"
	"github.com/Ankumeah/SourceFlameBackend/internal/external_auth"
	"github.com/Ankumeah/SourceFlameBackend/internal/session_store"

	"github.com/gin-gonic/gin"

	"net/http"
	"errors"
)

func login(r *gin.RouterGroup, app *a.App) {
  group := r.Group("/login")

  group.POST("", func(c *gin.Context) {
    ctx := c.Request.Context()
    var request struct {
      Username string `json:"username" binding:"required"`
      Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
      c.JSON(bad_request(err))
      return
    }

    user_id, err := app.User_db.Add_User(ctx, request.Username, request.Password)
    if errors.Is(err, database.Error_Exists) {
      user_id, err = app.User_db.Get_Id(ctx, request.Username)
      if errors.Is(err, database.Error_Invalid) {
        c.JSON(bad_request(err))
        return
      } else if err != nil {
        c.JSON(internal_server_error())
        return
      }

      valid, err := app.User_db.Verify_User(ctx, user_id, request.Password)
      if err != nil {
        c.JSON(internal_server_error())
        return
      } else if !valid {
        c.JSON(http.StatusUnauthorized, gin.H { "error": "Invalid password" })
        return
      }
    } else if err != nil {
      c.JSON(internal_server_error())
      return
    }

    add_session(c, app, request.Username)
  })

  group.POST("/external_login", func (c *gin.Context) {
    var request struct {
      JWT_type string `json:"JWT_type" binding:"required"`
      JWT string `json:"JWT" binding:"required"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
      c.JSON(bad_request(err))
      return
    }

    username, err := external_auth.Validate(request.JWT_type, request.JWT)
    if err == external_auth.Error_unsupported_JWT_type {
      c.JSON(http.StatusBadRequest, gin.H { "error": "Unsupported JWT type" })
      return
    } else if err == external_auth.Error_invalid_JWT {
      c.JSON(http.StatusUnauthorized, gin.H { "error": "Invalid JWT" })
      return
    } else if err != nil {
      c.JSON(internal_server_error())
      return
    }

    add_session(c, app, username)
  })
}

func add_session(c *gin.Context, app *a.App, username string) bool {
  ctx := c.Request.Context()

  token, err := app.Store.Add_Session(ctx, username)
  if err == session_store.Error_too_many_tokens {
    c.JSON(http.StatusForbidden, gin.H { "error": "Too many refresh tokens" })
    return false
  } else if err != nil {
    c.JSON(internal_server_error())
    return false
  }

  new_jwt, err := app.JWT_Handler.Issue_jwt(username)
  if err != nil {
    app.Store.Delete_Session(ctx, username, token)
    c.JSON(internal_server_error())
    return false
  }

  c.JSON(http.StatusOK, gin.H { "JWT": new_jwt, "refresh_token": token })
  return true
}
