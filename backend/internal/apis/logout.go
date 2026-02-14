package apis

import (
	"github.com/Ankumeah/DeltaBase/internal/session_store"

	"github.com/gin-gonic/gin"

	"net/http"
)

func logout(r *gin.RouterGroup, store *session_store.Session_store) {
  r.POST("/logout", func (c *gin.Context) {
    token, ok := validate_session(c, store)
    if !ok { return }

    delete_session(c, store, token)
  })
}

func validate_session(c *gin.Context, strore *session_store.Session_store) (string, bool) {
  ctx := c.Request.Context()

  var request struct {
    Refresh_token string `json:"refresh_token" binding:"required"`
    Email string `json:"email" binding:"required"`
  }

  if err := c.ShouldBindJSON(&request); err != nil {
    c.JSON(http.StatusBadRequest, gin.H { "error": err.Error() })
    return "", false
  }

  valid, err := strore.Validate_Session(ctx, request.Email, request.Refresh_token)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server error" })
    return "", false
  } else if !valid {
    c.JSON(http.StatusUnauthorized, gin.H { "error": "Invalid refresh token" })
    return "", false
  }

  return request.Refresh_token, true
}

func delete_session(c *gin.Context, store *session_store.Session_store, token string) bool {
  ctx := c.Request.Context()

  if err := store.Delete_Session(ctx, token); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server erro" })
    return  false
  }

  return  true
}
