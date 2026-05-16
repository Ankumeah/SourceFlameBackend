package apis

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/database"

	"github.com/gin-gonic/gin"

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
        get_repos(c, app)
      default:
        c.Status(http.StatusNotFound)
    }
  })
}

func meta(c *gin.Context, app *a.App) {
  ctx := c.Request.Context()
  username := c.Param("username")

  user_id, ok := get_user_id(c, app.User_db, username)
  if !ok { return }

  info, err := app.User_db.Info(ctx, user_id)
  if err != nil {
    c.JSON(internal_server_error())
    return
  }

  c.JSON(http.StatusOK, info)
}

func get_repos(c *gin.Context, app *a.App) {
  ctx := c.Request.Context()
  username := c.Param("username")

  _limit, err := strconv.ParseUint(c.Query("limit"), 10, 8)
  limit := uint8(_limit)
  if limit <= 0 { limit = 10 }
  if err != nil {
    c.JSON(bad_request(err))
    return
  }

  offset, err := strconv.ParseUint(c.Query("offset"), 10, 64)
  if err != nil {
    c.JSON(bad_request(err))
    return
  }

  user_id, ok := get_user_id(c, app.User_db, username)
  if !ok { return }

  repos, err := app.Git_db.Get_Repos(ctx, user_id, false, limit, offset)
  if err == database.Error_limit_too_big {
    c.JSON(http.StatusBadRequest, gin.H { "error": "Limit too big" })
    return
  } else if err != nil {
    c.JSON(internal_server_error())
    return
  }

  c.JSON(http.StatusOK, gin.H { "repos": repos })
}
