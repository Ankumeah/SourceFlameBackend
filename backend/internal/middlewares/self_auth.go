package middlewares

import (
	"github.com/Ankumeah/SourceFlameBackend/internal/database"

	"github.com/gin-gonic/gin"

	"errors"
	"net/http"
	"strings"
)

func SelfAuthMiddleware(userDb database.UserDb, create bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		parts := strings.Split(c.GetHeader("X-Authorization"), " ")
		if len(parts) != 2 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Incorrect segments in auth header"})
			return
		} else if parts[0] != selfAuthType {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Incorrect auth type, expected Basic"})
			return
		}

		username, password, err := parseBasicAuth(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(badHttpRequest(err))
			return
		}

		if create {
			userId, err := userDb.AddUser(ctx, username, password)
			if err == nil {
				c.Set(UserIdField, userId)
				c.Set(UsernameField, username)
				c.Next()
				return
			} else if !errors.Is(err, database.ErrExists) {
				c.AbortWithStatusJSON(internalServerError())
				return
			}
		}

		userId, err := userDb.GetId(ctx, username)
		if errors.Is(err, database.ErrInvalidUser) {
			c.AbortWithStatusJSON(badHttpRequest(err))
			return
		} else if err != nil {
			c.AbortWithStatusJSON(internalServerError())
			return
		}

		valid, err := userDb.VerifyUser(ctx, userId, password)
		if errors.Is(err, database.ErrInvalid) {
			c.AbortWithStatusJSON(badHttpRequest(err))
			return
		} else if err != nil {
			c.AbortWithStatusJSON(internalServerError())
			return
		} else if !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
			return
		} else {
			c.Set(UserIdField, userId)
			c.Set(UsernameField, username)
			c.Next()
			return
		}
	}
}
