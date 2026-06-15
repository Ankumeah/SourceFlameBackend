package apis

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/database"
	"github.com/Ankumeah/SourceFlameBackend/internal/git"
	"github.com/Ankumeah/SourceFlameBackend/internal/middlewares"

	"github.com/gin-gonic/gin"

	"errors"
	"fmt"
	"net/http"
	"strings"
)

func repo(r *gin.RouterGroup, app *a.App) {
  group := r.Group("/repo",
    middlewars.Check_Login_Middleware(app.JWT_Handler),
    middlewars.Check_PAT_Middleware(app),
  )

  group.GET("/:repo_owner/:repo_name/meta", func(c *gin.Context) {
    ctx := c.Request.Context()
    repo_name := c.Param("repo_name")
    repo_owner := c.Param("repo_owner")
    user_id := c.GetUint64("user_id")
    authed := c.GetBool("authed")

    owner_id, ok := get_user_id(c, app.User_db, repo_owner)
    if !ok { return }

    repo_id, err := app.Git_db.Get_Id(ctx, owner_id, repo_name)
    if err == database.Error_invalid_repo {
      c.JSON(invalid_repo())
      return
    } else if err != nil {
      c.JSON(internal_server_error())
      return
    }

    info, err := app.Git_db.Info(ctx, repo_id)
    if err != nil {
      c.JSON(internal_server_error())
      return
    }

    if info.Private && !authed {
      c.Header("WWW-Authenticate", app.Settings.AUTH_CHALLENGE_HEADER)
      c.JSON(unauthorised_request())
      return
    } else if info.Private && user_id != owner_id {
      c.JSON(bad_request(database.Error_invalid_repo))
      return
    }

    c.JSON(http.StatusOK, info)
  })

  group.GET("/:repo_owner/:repo_name/info/refs", func(c *gin.Context) {
    ctx := c.Request.Context()
    service := c.Query("service")
    repo_owner := c.Param("repo_owner")
    repo_name := c.Param("repo_name")
    user_id := c.GetUint64("user_id")
    authed := c.GetBool("authed")

    if service != "git-receive-pack" && service != "git-upload-pack" {
      c.JSON(http.StatusForbidden, gin.H { "error": "Unsupported Service" })
      return
    }

    owner_id, ok := get_user_id(c, app.User_db, repo_owner)
    if !ok { return }

    repo_id, err := app.Git_db.Get_Id(ctx, owner_id, repo_name)
    if err == database.Error_invalid_repo {
      c.JSON(invalid_repo())
      return
    } else if err != nil {
      c.JSON(internal_server_error())
      return
    }

    info, err := app.Git_db.Info(ctx, repo_id)
    if err != nil {
      c.JSON(internal_server_error())
      return
    }

    if info.Private && !authed {
      c.Header("WWW-Authenticate", app.Settings.AUTH_CHALLENGE_HEADER)
      c.JSON(unauthorised_request())
      return
    } else if info.Private && user_id != owner_id {
      c.JSON(bad_request(database.Error_invalid_repo))
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

  group.POST("/:repo_owner/:repo_name/git-upload-pack", func(c *gin.Context) {
    ctx := c.Request.Context()
    repo_owner := c.Param("repo_owner")
    repo_name := c.Param("repo_name")
    user_id := c.GetUint64("user_id")
    authed := c.GetBool("authed")

    owner_id, ok := get_user_id(c, app.User_db, repo_owner)
    if !ok { return }

    repo_id, err := app.Git_db.Get_Id(ctx, owner_id, repo_name)
    if err == database.Error_invalid_repo {
      c.JSON(invalid_repo())
      return
    } else if err != nil {
      c.JSON(internal_server_error())
      return
    }

    info, err := app.Git_db.Info(ctx, repo_id)
    if err != nil {
      c.JSON(internal_server_error())
      return
    }

    if info.Private && !authed {
      c.Header("WWW-Authenticate", app.Settings.AUTH_CHALLENGE_HEADER)
      c.JSON(unauthorised_request())
      return
    } else if info.Private && user_id != owner_id {
      c.JSON(bad_request(database.Error_invalid_repo))
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

  group.POST("/:repo_owner/:repo_name/git-receive-pack", func(c *gin.Context) {
    ctx := c.Request.Context()
    repo_owner := c.Param("repo_owner")
    repo_name := c.Param("repo_name")
    user_id := c.GetUint64("user_id")
    authed := c.GetBool("authed")

    owner_id, ok := get_user_id(c, app.User_db, repo_owner)
    if !ok { return }

    repo_id, err := app.Git_db.Get_Id(ctx, owner_id, repo_name)
    if err == database.Error_invalid_repo {
      c.JSON(invalid_repo())
      return
    } else if err != nil {
      c.JSON(internal_server_error())
      return
    }

    if !authed {
      c.Header("WWW-Authenticate", app.Settings.AUTH_CHALLENGE_HEADER)
      c.JSON(unauthorised_request())
      return
    } else if user_id != owner_id {
      c.JSON(bad_request(database.Error_invalid_repo))
      return
    }

    err = git.Receive_Pack(ctx, repo_id, c.Request.Body, c.Writer)
    if err != nil {
      c.JSON(internal_server_error())
      return
    }
  })

  group.GET("/:repo_owner/:repo_name/blob/*path", func(c *gin.Context) {
    ctx := c.Request.Context()
    repo_owner := c.Param("repo_owner")
    repo_name := c.Param("repo_name")
    user_id := c.GetUint64("user_id")
    authed := c.GetBool("authed")
    path := strings.Trim(c.Param("path"), "/")
    hash := c.Query("hash")

    owner_id, ok := get_user_id(c, app.User_db, repo_owner)
    if !ok { return }

    repo_id, err := app.Git_db.Get_Id(ctx, owner_id, repo_name)
    if err == database.Error_invalid_repo {
      c.JSON(invalid_repo())
      return
    } else if err != nil {
      c.JSON(internal_server_error())
      return
    }

    info, err := app.Git_db.Info(ctx, repo_id)
    if err != nil {
      c.JSON(internal_server_error())
      return
    }

    if info.Private && !authed {
      c.Header("WWW-Authenticate", app.Settings.AUTH_CHALLENGE_HEADER)
      c.JSON(unauthorised_request())
      return
    } else if info.Private && user_id != owner_id {
      c.JSON(bad_request(database.Error_invalid_repo))
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

  group.GET("/:repo_owner/:repo_name/list/*path", func(c *gin.Context) {
    ctx := c.Request.Context()
    repo_owner := c.Param("repo_owner")
    repo_name := c.Param("repo_name")
    user_id := c.GetUint64("user_id")
    authed := c.GetBool("authed")
    path := strings.Trim(c.Param("path"), "/")
    hash := c.Query("hash")

    owner_id, ok := get_user_id(c, app.User_db, repo_owner)
    if !ok { return }

    repo_id, err := app.Git_db.Get_Id(ctx, owner_id, repo_name)
    if err == database.Error_invalid_repo {
      c.JSON(invalid_repo())
      return
    } else if err != nil {
      c.JSON(internal_server_error())
      return
    }

    info, err := app.Git_db.Info(ctx, repo_id)
    if err != nil {
      c.JSON(internal_server_error())
      return
    }

    if info.Private && !authed {
      c.Header("WWW-Authenticate", app.Settings.AUTH_CHALLENGE_HEADER)
      c.JSON(unauthorised_request())
      return
    } else if info.Private && user_id != owner_id {
      c.JSON(bad_request(database.Error_invalid_repo))
      return
    }

    files, err := git.List_Dir(repo_id, hash, path)
    if errors.Is(err, git.Error_Not_Found) {
      c.JSON(http.StatusNotFound, gin.H { "error": err.Error() })
      return
    } else if err == git.Error_Path_Too_Deep {
      c.JSON(http.StatusRequestURITooLong, gin.H { "error": err.Error() })
      return
    } else if err != nil {
      c.JSON(internal_server_error())
      return
    }

    c.JSON(http.StatusOK, gin.H { "files": files })
  })
}
