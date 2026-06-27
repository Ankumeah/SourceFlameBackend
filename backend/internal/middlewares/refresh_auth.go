package middlewars

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"

	"github.com/gin-gonic/gin"

	"log"
	"net/http"
	"strings"
)

func Refresh_Auth_Middleware(app *a.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		username := c.GetHeader("X-Username")

		segments := strings.Split(c.GetHeader("X-Authorization"), " ")
		if len(segments) != 2 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Incorrect number of segments in auth header"})
			return
		}

		if segments[0] != refresh_auth_type {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Incorrect auth type"})
			return
		}

		valid, err := app.Store.Validate_Session(ctx, username, segments[1])
		if err != nil {
			log.Printf("Error while validateing session: %v\n", err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		} else if !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			return
		}

		c.Set(Username_feild, username)
		c.Set(Session_feild, segments[1])
		c.Next()
	}
}
