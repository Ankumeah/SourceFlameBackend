package middlewars

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/database"
	"github.com/Ankumeah/SourceFlameBackend/internal/jwt"

	"github.com/gin-gonic/gin"

	"net/http"
	"strings"
)

func JWT_Auth_Middleware(app *a.App, enforce bool) gin.HandlerFunc {
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

		if segments[0] != jwt_auth_type {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Incorrect auth type"})
			return
		}

		username, err := app.JWT_Handler.Validate_jwt(segments[1])
		if err == jwt.Error_invalid_JWT {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT"})
			return
		} else if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user_id, err := app.User_db.Get_Id(ctx, username)
		if err == database.Error_invalid_user {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid user"})
			return
		} else if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.Set(User_id_feild, user_id)
		c.Next()
	}
}
