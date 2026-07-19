package apis

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/database"

	"github.com/gin-gonic/gin"

	"errors"
	"net/http"
	"strconv"
)

func user(r *gin.RouterGroup, app *a.App) {
	group := r.Group("/user")

	group.GET("/:username/*action", func(c *gin.Context) {
		switch c.Param("action") {
		case "/meta":
			meta(c, app)
		case "/repos":
			getRepos(c, app)
		default:
			c.Status(http.StatusNotFound)
		}
	})
}

func meta(c *gin.Context, app *a.App) {
	ctx := c.Request.Context()
	username := c.Param("username")

	userId, ok := getUserId(c, app.UserDb, username)
	if !ok {
		return
	}

	info, err := app.UserDb.Info(ctx, userId)
	if err != nil {
		c.JSON(internalServerError())
		return
	}

	c.JSON(http.StatusOK, info)
}

func getRepos(c *gin.Context, app *a.App) {
	ctx := c.Request.Context()
	username := c.Param("username")

	_limit, err := strconv.ParseUint(c.Query("limit"), 10, 8)
	limit := uint8(_limit)
	if limit <= 0 {
		limit = 10
	}
	if err != nil {
		c.JSON(badRequest(err))
		return
	}

	offset, err := strconv.ParseUint(c.Query("offset"), 10, 64)
	if err != nil {
		c.JSON(badRequest(err))
		return
	}

	userId, ok := getUserId(c, app.UserDb, username)
	if !ok {
		return
	}

	repos, err := app.GitDb.GetRepos(ctx, userId, limit, offset)
	if errors.Is(err, database.ErrLimitTooLarge) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Limit too big"})
		return
	} else if err != nil {
		c.JSON(internalServerError())
		return
	}

	c.JSON(http.StatusOK, gin.H{"repos": repos})
}
