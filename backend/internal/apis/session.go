package apis

import (
	"errors"

	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/middlewares"
	"github.com/Ankumeah/SourceFlameBackend/internal/session_store"

	"github.com/gin-gonic/gin"

	"net/http"
)

func session(r *gin.RouterGroup, app *a.App) {
	group := r.Group("/session", middlewares.RefreshAuthMiddleware(app))

	group.POST("/renew_jwt", func(c *gin.Context) {
		username := c.GetString(middlewares.UsernameField)

		token, err := app.JWTHandler.IssueJwt(username)
		if err != nil {
			c.JSON(internalServerError())
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"JWT": token})
			return
		}
	})

	group.POST("/renew_session", func(c *gin.Context) {
		ctx := c.Request.Context()
		username := c.GetString(middlewares.UsernameField)
		session := c.GetString(middlewares.SessionField)

		token, err := app.Store.AddSession(ctx, username)
		if errors.Is(err, session_store.ErrTooManyTokens) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Too many tokens"})
			return
		} else if err != nil {
			c.JSON(internalServerError())
			return
		}

		if err = app.Store.DeleteSession(ctx, username, session); err != nil {
			c.JSON(internalServerError())
			app.Store.DeleteSession(ctx, username, token)
			return
		}

		c.JSON(http.StatusOK, gin.H{"refreshToken": token})
	})

	group.DELETE("", func(c *gin.Context) {
		ctx := c.Request.Context()
		session := c.GetString(middlewares.SessionField)
		username := c.GetString(middlewares.UsernameField)

		if err := app.Store.DeleteSession(ctx, username, session); err != nil {
			c.JSON(internalServerError())
			return
		}

		c.Status(http.StatusOK)
	})
}
