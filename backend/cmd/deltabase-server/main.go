package main

import (
	"github.com/Ankumeah/DeltaBase/internal/apis"

	"github.com/gin-gonic/gin"

	"fmt"
  "os"
)

func main() {
	fmt.Println("Starting http server")

	r := gin.Default()

	apiGroup := r.Group("/api")
	apis.Apis(apiGroup)

  port, ok := os.LookupEnv("PORT")
  if !ok {
    panic("env var PORT not set")
  }

  fmt.Println("Running backend on port: " + port)
	r.Run(":" + port)
}
