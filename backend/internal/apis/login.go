package apis

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/middlewares"
	"github.com/Ankumeah/SourceFlameBackend/internal/session_store"

	"github.com/gin-gonic/gin"

	"errors"
	"net/http"
)

func login(r *gin.RouterGroup, app *a.App) {
	group := r.Group("/login", middlewares.SelfAuthMiddleware(*app.UserDb, true))

	group.POST("", func(c *gin.Context) {
		ctx := c.Request.Context()
		username := c.GetString(middlewares.UsernameField)

		token, err := app.Store.AddSession(ctx, username)
		if errors.Is(err, session_store.ErrTooManyTokens) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Too many refresh tokens"})
			return
		} else if err != nil {
			c.JSON(internalServerError())
			return
		}

		newJwt, err := app.JWTHandler.IssueJwt(username)
		if err != nil {
			app.Store.DeleteSession(ctx, username, token)
			c.JSON(internalServerError())
			return
		}

		c.JSON(http.StatusOK, gin.H{"JWT": newJwt, "refresh_token": token})
	})
}
