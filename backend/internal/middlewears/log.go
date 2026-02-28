package middlewears

import (
	"github.com/gin-gonic/gin"

  "log"
)

func Log_Middlewear() gin.HandlerFunc {
  return func(c *gin.Context) {
    ip := c.ClientIP()
    route := c.FullPath()

    log.Printf("Got request: %v | %v\n", ip, route)

    c.Next()
  }
}
