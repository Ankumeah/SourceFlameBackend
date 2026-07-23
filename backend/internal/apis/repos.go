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
	g := r.Group("/repos", middlewares.JWTAuthMiddleware(app, true))
  group := g.Group("/:repo_name")

	group.POST("", func(c *gin.Context) {
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

	group.DELETE("", func(c *gin.Context) {
		ctx := c.Request.Context()
		userId := c.GetUint64(middlewares.UserIdField)
		repoName := c.Param("repo_name")

    repoId, ok := getRepoId(c, app.GitDb, userId, repoName)
    if !ok {
      return
    }

    err := app.GitDb.DeleteRepo(ctx, repoId)
		if errors.Is(err, database.ErrInvalid) {
			c.JSON(badRequest(err))
			return
		} else if err != nil {
			c.JSON(internalServerError())
			return
		}

		c.Status(http.StatusOK)
	})

  group.POST("/transfer/:new_owner", func(c *gin.Context) {
    ctx := c.Request.Context()
    userId := c.GetUint64(middlewares.UserIdField)
    repoName := c.Param("repo_name")
    newOwner := c.Param("new_owner")

    newOwnerId, ok := getUserId(c, app.UserDb, newOwner)
    if !ok {
      return
    }
    repoId, ok := getRepoId(c, app.GitDb, userId, repoName)
    if !ok {
      return
    }

    err := app.GitDb.TransferOwner(ctx, repoId, newOwnerId)
    if errors.Is(err, database.ErrSafe) {
      c.JSON(badRequest(err))
    } else if err != nil {
      c.JSON(internalServerError())
      return
    }

    c.Status(http.StatusOK)
  })
}
