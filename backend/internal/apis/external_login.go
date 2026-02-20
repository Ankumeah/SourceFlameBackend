package apis

import (
	"github.com/Ankumeah/DeltaBase/internal/external_auth"
	"github.com/Ankumeah/DeltaBase/internal/jwt"
	"github.com/Ankumeah/DeltaBase/internal/session_store"

	"github.com/gin-gonic/gin"

	"net/http"
)

func external_login(r *gin.RouterGroup, store *session_store.Session_store) {
  r.POST("/external_login", func (c *gin.Context) {
    username, ok := validate_JWT(c)
    if !ok { return }

    add_session(c, store, username)
  })
}

func validate_JWT(c *gin.Context) (string, bool) {
  var request struct {
    JWT_type string `json:"JWT_type" binding:"required"`
    JWT string `json:"JWT" binding:"required"`
  }

  if err := c.ShouldBindJSON(&request); err != nil {
    c.JSON(http.StatusBadRequest, gin.H { "error": err.Error() })
    return "", false
  }

  username, err := external_auth.Validate(request.JWT_type, request.JWT)
  if err == external_auth.Error_unsupported_JWT_type {
    c.JSON(http.StatusBadRequest, gin.H { "error": "Unsupported JWT type" })
    return "", false
  } else if err == external_auth.Error_invalid_JWT {
    c.JSON(http.StatusUnauthorized, gin.H { "error": "Invalid JWT" })
    return "", false
  } else if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server error" })
    return "", false
  }

  return username, true
}

func add_session(c *gin.Context, store *session_store.Session_store, username string) bool {
  ctx := c.Request.Context()

  token, err := store.Add_Session(ctx, username)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H { "error": "Invalid server error" })
    return false
  }

  new_jwt, err := jwt.Issue_jwt(username)
  if err != nil {
    store.Delete_Session(ctx, token)

    c.JSON(http.StatusInternalServerError, gin.H { "error": "Invalid server error" })
    return false
  } else {
    c.JSON(http.StatusOK, gin.H { "JWT": new_jwt, "refresh_token": token })
    return  true
  }
}
