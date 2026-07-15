package apis

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/middlewares"
	"github.com/Ankumeah/SourceFlameBackend/internal/session_store"

	"github.com/gin-gonic/gin"

	"net/http"
)

func login(r *gin.RouterGroup, app *a.App) {
	group := r.Group("/login", middlewars.Self_Auth_Middleware(*app.User_db, true))

	group.POST("", func(c *gin.Context) {
		ctx := c.Request.Context()
		username := c.GetString(middlewars.Username_feild)

		token, err := app.Store.Add_Session(ctx, username)
		if err == session_store.Error_too_many_tokens {
			c.JSON(http.StatusForbidden, gin.H{"error": "Too many refresh tokens"})
			return
		} else if err != nil {
			c.JSON(internal_server_error())
			return
		}

		new_jwt, err := app.JWT_Handler.Issue_jwt(username)
		if err != nil {
			app.Store.Delete_Session(ctx, username, token)
			c.JSON(internal_server_error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"JWT": new_jwt, "refresh_token": token})
	})
}
