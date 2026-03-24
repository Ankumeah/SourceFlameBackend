package apis

import (
	"github.com/Ankumeah/DeltaBase/internal/database"
	"github.com/Ankumeah/DeltaBase/internal/middlewears"

	"github.com/gin-gonic/gin"

	"net/http"
  "fmt"
)

func repos(r *gin.RouterGroup, db *database.Git_db) {
  group := r.Group("/repos", middlewears.Verify_JWT_Middlewear())

  group.POST("/create", func (c *gin.Context) {
    ctx := c.Request.Context()
    _username, _ := c.Get("username")
    username := fmt.Sprintf("%v", _username)
    var request struct {
      Repo_name string `json:"repo_name" binding:"required"`
      Private bool `json:"private" binding:"required"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
      c.JSON(http.StatusBadRequest, gin.H { "error": err.Error() })
      return
    }

    _, err := db.Get_Id(ctx, username, request.Repo_name)
    if err == database.Error_invalid_user {
      c.JSON(http.StatusBadRequest, gin.H { "error": "Invalid user" })
      return
    } else if err != database.Error_invalid_repo && err != nil {
      c.JSON(http.StatusInternalServerError, gin.H { "error": " Internal server error" })
      return
    } else if err != database.Error_invalid_repo {
      c.JSON(http.StatusConflict, gin.H { "error": "Repo alreday exists" })
      return
    }

    repo_id, err := db.Create_Repo(ctx, username, request.Repo_name, request.Private)
    if err == database.Error_invalid_user {
      c.JSON(http.StatusBadRequest, gin.H { "error": "Invalid user" })
      return
    } else if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server error" })
      return
    } else {
      c.JSON(http.StatusOK, gin.H { "repo_id": repo_id })
    }
  })
}
