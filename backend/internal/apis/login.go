package apis

import (
	"github.com/Ankumeah/DeltaBase/internal/auth/jwt"

	"github.com/gin-gonic/gin"

	"net/http"
)

func login(r *gin.RouterGroup) {
  r.POST("/login", handle_login)
}

var supported_JWT_types = map[string]jwt.Handler {
  "google": jwt.Google,
}

func handle_login(c *gin.Context) {
  var request struct {
    JWT_type string `json:"JWT_type" binding:"required"`
    JWT string `json:"JWT" binding:"required"`
  }

  if err := c.ShouldBindJSON(&request); err != nil {
    c.JSON(http.StatusBadRequest, gin.H { "error": err.Error() })
    return
  }

  handler, ok := supported_JWT_types[request.JWT_type]
  if ok != true {
    c.JSON(http.StatusBadRequest, gin.H { "error": "Unsupported JWT type" })
    return
  }

  res, err := handler(request.JWT)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server error" })
    return
  }

  if res != true {
    c.JSON(http.StatusUnauthorized, gin.H { "error": "Invalid JWT" })
    return
  }

  // TODO("Add login")

  c.JSON(http.StatusOK, gin.H {})
}
