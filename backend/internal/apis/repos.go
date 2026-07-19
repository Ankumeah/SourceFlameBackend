package apis

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/database"
	"github.com/Ankumeah/SourceFlameBackend/internal/middlewares"

	"github.com/gin-gonic/gin"

	"errors"
	"net/http"
)

func repos(r *gin.RouterGroup, app *a.App) {
	group := r.Group("/repos", middlewares.JWTAuthMiddleware(app, true))

	group.POST("/:repo_name", func(c *gin.Context) {
		ctx := c.Request.Context()
		repoName := c.Param("repo_name")
		userId := c.GetUint64(middlewares.UserIdField)

		_, err := app.GitDb.CreateRepo(ctx, userId, repoName)
		if errors.Is(err, database.ErrSafe) {
			c.JSON(badRequest(err))
			return
		} else if err != nil {
			c.JSON(internalServerError())
			return
		}

		c.Status(http.StatusOK)
	})

	group.DELETE("/:repo_name", func(c *gin.Context) {
		ctx := c.Request.Context()
		userId := c.GetUint64(middlewares.UserIdField)
		repoName := c.Param("repo_name")

		repoId, err := app.GitDb.GetId(ctx, userId, repoName)
		if errors.Is(err, database.ErrInvalid) {
			c.JSON(badRequest(err))
			return
		} else if err != nil {
			c.JSON(internalServerError())
			return
		}

		err = app.GitDb.DeleteRepo(ctx, repoId)
		if errors.Is(err, database.ErrInvalid) {
			c.JSON(badRequest(err))
			return
		} else if err != nil {
			c.JSON(internalServerError())
			return
		}

		c.Status(http.StatusOK)
	})
}
