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
	group := r.Group("/repos", middlewars.JWT_Auth_Middleware(app, true))

	group.POST("/:repo_name", func(c *gin.Context) {
		ctx := c.Request.Context()
		repo_name := c.Param("repo_name")
		user_id := c.GetUint64("user_id")

		_, err := app.Git_db.Create_Repo(ctx, user_id, repo_name)
		if errors.Is(err, database.Safe_Error) {
			c.JSON(bad_request(err))
			return
		} else if err != nil {
			c.JSON(internal_server_error())
			return
		}

		c.Status(http.StatusOK)
	})

	group.DELETE("/:repo_name", func(c *gin.Context) {
		ctx := c.Request.Context()
		user_id := c.GetUint64("user_id")
		repo_name := c.Param("repo_name")

		repo_id, err := app.Git_db.Get_Id(ctx, user_id, repo_name)
		if errors.Is(err, database.Error_Invalid) {
			c.JSON(invalid_user())
			return
		} else if err != nil {
			c.JSON(internal_server_error())
			return
		}

		err = app.Git_db.Delete_Repo(ctx, repo_id)
		if errors.Is(err, database.Error_Invalid) {
			c.JSON(bad_request(err))
			return
		} else if err != nil {
			c.JSON(internal_server_error())
			return
		}

		c.Status(http.StatusOK)
	})
}
