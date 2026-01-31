package apis

import (
  "github.com/gin-gonic/gin"
)

func Apis(r *gin.RouterGroup) {
  login(r)
}
