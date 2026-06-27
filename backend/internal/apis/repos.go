package apis

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/database"
	"github.com/Ankumeah/SourceFlameBackend/internal/middlewares"

	"github.com/gin-gonic/gin"

	"errors"
	"net/http"
	"strconv"
)

func repos(r *gin.RouterGroup, app *a.App) {
	group := r.Group("/repos", middlewars.JWT_Auth_Middleware(app, true))

	group.POST("/:repo_name", func(c *gin.Context) {
		ctx := c.Request.Context()
		repo_name := c.Param("repo_name")
		user_id := c.GetUint64("user_id")

		private, err := strconv.ParseBool(c.Query("private"))
		if err != nil {
			c.JSON(bad_request(err))
			return
		}

		_, err = app.Git_db.Create_Repo(ctx, user_id, repo_name, private)
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

	group.GET("/all", func(c *gin.Context) {
		ctx := c.Request.Context()
		user_id := c.GetUint64("user_id")

		_limit, err := strconv.ParseUint(c.Query("limit"), 10, 8)
		limit := uint8(_limit)
		if limit <= 0 {
			limit = 10
		}
		if err != nil {
			c.JSON(bad_request(err))
			return
		}

		offset, err := strconv.ParseUint(c.Query("offset"), 10, 64)
		if err != nil {
			c.JSON(bad_request(err))
			return
		}

		repos, err := app.Git_db.Get_Repos(ctx, user_id, true, limit, offset)
		if err == database.Error_limit_too_big {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Limit too big"})
			return
		} else if err != nil {
			c.JSON(internal_server_error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"repos": repos})
	})
}
