package apis

import (
	"github.com/Ankumeah/DeltaBase/internal/jwt"

	"github.com/gin-gonic/gin"

	"net/http"
  "fmt"
)

func renew_JWT(r *gin.RouterGroup) {
  r.POST("/renew_jwt", func (c *gin.Context) {
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
}
