package apis

import (
	"github.com/Ankumeah/DeltaBase/internal/auth/session"

	"github.com/gin-gonic/gin"

	"net/http"
)

func renew_JWT(r *gin.RouterGroup) {
  r.POST("/renew_jwt", func (c *gin.Context) {
    username := c.GetHeader("username")

    token, err := session.Issue_jwt(username)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server error" })
      return
    } else {
      c.JSON(http.StatusOK, gin.H { "JWT": token })
      return
    }
  })
}
