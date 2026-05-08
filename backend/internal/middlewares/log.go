package middlewars

import (
	"github.com/gin-gonic/gin"

  "log"
)

func Log_Middleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    log.Printf("New request from agent: %v\n", c.GetHeader("User-Agent"))

    c.Next()
  }
}
