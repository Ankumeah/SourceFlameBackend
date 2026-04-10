package apis

import (
	"github.com/Ankumeah/DeltaBase/internal/database"
	"github.com/Ankumeah/DeltaBase/internal/git"

	"github.com/gin-gonic/gin"

	"net/http"
  "fmt"
  "strings"
  "errors"
)

func repo(r *gin.RouterGroup, git_db *database.Git_db, user_db *database.User_db) {
  group := r.Group("/repo")

  group.GET("/:username/:repo_name/meta", func(c *gin.Context) {
    ctx := c.Request.Context()
    repo_name := c.Param("repo_name")
    username := c.Param("username")

    owner_id, ok := get_user_id(c, user_db, username)
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

  group.GET("/:username/:repo_name/info/refs", func(c *gin.Context) {
    ctx := c.Request.Context()
    service := c.Query("service")
    username := c.Param("username")
    repo_name := c.Param("repo_name")

    if service != "git-receive-pack" && service != "git-upload-pack" {
      c.JSON(http.StatusForbidden, gin.H { "error": "Unsupported Service" })
      return
    }

    owner_id, ok := get_user_id(c, user_db, username)
    if !ok { return }

    repo_id, err := git_db.Get_Id(ctx, owner_id, repo_name)
    if err == database.Error_invalid_repo {
      c.JSON(invalid_repo())
      return
    } else if err != nil {
      c.JSON(internal_server_error())
      return
    }

    c.Header("Content-Type", fmt.Sprintf("application/x-%s-advertisement", service))
    c.Header("Cache-Control", "no-cache")

    err = git.Info_Refs(ctx, repo_id, service, c.Writer)
    if err != nil {
      c.JSON(internal_server_error())
      return
    }
  })

  group.POST("/:username/:repo_name/git-upload-pack", func(c *gin.Context) {
    ctx := c.Request.Context()
    username := c.Param("username")
    repo_name := c.Param("repo_name")

    owner_id, ok := get_user_id(c, user_db, username)
    if !ok { return }

    repo_id, err := git_db.Get_Id(ctx, owner_id, repo_name)
    if err == database.Error_invalid_repo {
      c.JSON(invalid_repo())
      return
    } else if err != nil {
      c.JSON(internal_server_error())
      return
    }

    c.Header("Content-Type", "application/x-git-upload-pack-result")
    c.Header("Cache-Control", "no-cache")

    err = git.Upload_Pack(ctx, repo_id, c.Request.Body, c.Writer)
    if err != nil {
      c.JSON(internal_server_error())
      return
    }
  })

  group.POST("/:username/:repo_name/git-receive-pack", func(c *gin.Context) {
    ctx := c.Request.Context()
    username := c.Param("username")
    repo_name := c.Param("repo_name")

    owner_id, ok := get_user_id(c, user_db, username)
    if !ok { return }

    repo_id, err := git_db.Get_Id(ctx, owner_id, repo_name)
    if err == database.Error_invalid_repo {
      c.JSON(invalid_repo())
      return
    } else if err != nil {
      c.JSON(internal_server_error())
      return
    }

    err = git.Receive_Pack(ctx, repo_id, c.Request.Body, c.Writer)
    if err != nil {
      c.JSON(internal_server_error())
      return
    }
  })

  group.GET("/:username/:repo_name/blob/*path", func(c *gin.Context) {
    ctx := c.Request.Context()
    username := c.Param("username")
    repo_name := c.Param("repo_name")
    path := strings.Trim(c.Param("path"), "/")
    hash := c.Query("hash")

    owner_id, ok := get_user_id(c, user_db, username)
    if !ok { return }

    repo_id, err := git_db.Get_Id(ctx, owner_id, repo_name)
    if err == database.Error_invalid_repo {
      c.JSON(invalid_repo())
      return
    } else if err != nil {
      c.JSON(internal_server_error())
      return
    }

    blob, err := git.Get_Glob(repo_id, hash, path)
    if errors.Is(err, git.Error_Not_Found) {
      c.JSON(http.StatusNotFound, gin.H { "error": err.Error() })
      return
    } else if err == git.Error_Blob_Too_Large {
      c.JSON(http.StatusRequestEntityTooLarge, gin.H { "error": err.Error() })
      return
    } else if err != nil {
      c.JSON(internal_server_error())
      return
    }

    blob_type := http.DetectContentType([]byte(blob))
    c.Header("Content-Type", blob_type)
    c.String(http.StatusOK, blob)
  })

  group.GET("/:username/:repo_name/list/*path", func(c *gin.Context) {
    ctx := c.Request.Context()
    username := c.Param("username")
    repo_name := c.Param("repo_name")
    path := strings.Trim(c.Param("path"), "/")
    hash := c.Query("hash")

    owner_id, ok := get_user_id(c, user_db, username)
    if !ok { return }

    repo_id, err := git_db.Get_Id(ctx, owner_id, repo_name)
    if err == database.Error_invalid_repo {
      c.JSON(invalid_repo())
      return
    } else if err != nil {
      c.JSON(internal_server_error())
      return
    }

    files, err := git.List_Dir(repo_id, hash, path)
    if err == git.Error_Inavlid_Commit_Hash ||
      err == git.Error_Blob_Not_Found ||
      err == git.Error_Path_Not_Found ||
      err == git.Error_Path_Too_Deep {

      c.JSON(http.StatusNotFound, gin.H { "error": err.Error() })
      return
    } else if err != nil {
      c.JSON(internal_server_error())
      return
    }

    c.JSON(http.StatusOK, files)
  })
}
