package middlewares

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/database"
	"github.com/Ankumeah/SourceFlameBackend/internal/jwt"

	"github.com/gin-gonic/gin"

	"net/http"
  "errors"
	"strings"
)

func JWTAuthMiddleware(app *a.App, enforce bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		token := c.GetHeader("X-Authorization")
		if token == "" {
			if enforce {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "JWT not provided"})
				return
			}
			c.Next()
			return
		}

		segments := strings.Split(token, " ")
		if len(segments) != 2 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Incorrect number of segments in auth header"})
			return
		}

		if segments[0] != jwtAuthType {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Incorrect auth type"})
			return
		}

		username, err := app.JWTHandler.ValidateJwt(segments[1])
		if errors.Is(err, jwt.ErrInvalidJWT) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT"})
			return
		} else if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userId, err := app.UserDb.GetId(ctx, username)
		if errors.Is(err, database.ErrInvalidUser) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid user"})
			return
		} else if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.Set(UserIdField, userId)
		c.Next()
	}
}
