package middlewears

import (
	"github.com/gin-gonic/gin"

  "log"
)

func Log_Middlewear() gin.HandlerFunc {
  return func(c *gin.Context) {
    ip := c.ClientIP()
    route := c.FullPath()

    log.Printf("Got request from %v at %v\n", ip, route)

    c.Next()

    status := c.Writer.Status()

    log.Printf("Closed request from %v at %v with %v\n", ip, route, status)
  }
}
