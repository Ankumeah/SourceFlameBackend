package middlewars

import (
	"github.com/Ankumeah/SourceFlameBackend/internal/database"

	"github.com/gin-gonic/gin"

	"errors"
	"net/http"
	"strings"
)

func Self_Auth_Middleware(user_db database.User_db, create bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		parts := strings.Split(c.GetHeader("X-Authorization"), " ")
		if len(parts) != 2 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Incorrect segments in auth header"})
			return
		} else if parts[0] != self_auth_type {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Incorrect auth type, expected Basic"})
			return
		}

		username, password, err := parse_basic_auth(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(bad_http_request(err))
			return
		}

		if create {
			user_id, err := user_db.Add_User(ctx, username, password)
			if err == nil {
				c.Set(User_id_feild, user_id)
				c.Set(Username_feild, username)
				c.Next()
				return
			} else if !errors.Is(err, database.Error_Exists) {
				c.AbortWithStatusJSON(internal_server_error())
				return
			}
		}

		user_id, err := user_db.Get_Id(ctx, username)
		if err == database.Error_invalid_user {
			c.AbortWithStatusJSON(bad_http_request(err))
			return
		} else if err != nil {
			c.AbortWithStatusJSON(internal_server_error())
			return
		}

		valid, err := user_db.Verify_User(ctx, user_id, password)
		if errors.Is(err, database.Error_Invalid) {
			c.AbortWithStatusJSON(bad_http_request(err))
			return
		} else if err != nil {
			c.AbortWithStatusJSON(internal_server_error())
			return
		} else if !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
			return
		} else {
			c.Set(User_id_feild, user_id)
			c.Set(Username_feild, username)
			c.Next()
			return
		}
	}
}
