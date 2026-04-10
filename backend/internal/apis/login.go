package apis

import (
	"github.com/Ankumeah/DeltaBase/internal/database"
	"github.com/Ankumeah/DeltaBase/internal/jwt"
	"github.com/Ankumeah/DeltaBase/internal/session_store"
	"github.com/Ankumeah/DeltaBase/internal/external_auth"
	a "github.com/Ankumeah/DeltaBase/internal/app"

	"github.com/gin-gonic/gin"

	"net/http"
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

    _, err := app.User_db.Get_Id(ctx, request.Username)
    if err == database.Error_invalid_user {
      _, err = app.User_db.Add_User(ctx, request.Username, request.Password)
    }
    if err != nil {
      c.JSON(internal_server_error())
      return
    }

    add_session(c, app.Store, request.Username)
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

    add_session(c, app.Store, username)
  })
}

func add_session(c *gin.Context, store *session_store.Session_store, username string) bool {
  ctx := c.Request.Context()

  token, err := store.Add_Session(ctx, username)
  if err == session_store.Error_too_many_tokens {
    c.JSON(http.StatusForbidden, gin.H { "error": "Too many refresh tokens" })
    return false
  } else if err != nil {
    c.JSON(internal_server_error())
    return false
  }

  new_jwt, err := jwt.Issue_jwt(username)
  if err != nil {
    store.Delete_Session(ctx, username, token)

    c.JSON(internal_server_error())
    return false
  } else {
    c.JSON(http.StatusOK, gin.H { "JWT": new_jwt, "refresh_token": token })
    return true
  }
}
