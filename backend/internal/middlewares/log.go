package middlewars

import (
	"github.com/gin-gonic/gin"
)

func Log_Middleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    c.Next()
  }
}
