package apis

import (
	"github.com/Ankumeah/DeltaBase/internal/database"

	"github.com/gin-gonic/gin"

	"net/http"
)

func repo(r *gin.RouterGroup, git_db *database.Git_db, user_db *database.User_db) {
  group := r.Group("/repo")

  group.GET("/:username/:repo_name/meta", func(c *gin.Context) {
    ctx := c.Request.Context()
    repo_name := c.Param("repo_name")

    owner_id, ok := get_id(c, user_db)
    if !ok { return }

    repo_id, err := git_db.Get_Id(ctx, owner_id, repo_name)
    if err == database.Error_invalid_repo {
      c.JSON(invalid_repo())
      return
    } else if err != nil {
      c.JSON(internal_server_error())
      return
    }

    info, err := git_db.Info(ctx, repo_id)
    if err != nil {
      c.JSON(internal_server_error())
      return
    }

    c.JSON(http.StatusOK, info)
  })
}
