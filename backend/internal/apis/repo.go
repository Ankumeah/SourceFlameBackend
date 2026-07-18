package apis

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/database"
	"github.com/Ankumeah/SourceFlameBackend/internal/git"
	"github.com/Ankumeah/SourceFlameBackend/internal/middlewares"
	"github.com/Ankumeah/SourceFlameBackend/internal/roles"

	"github.com/gin-gonic/gin"

	"errors"
	"fmt"
	"net/http"
	"strings"
)

func repo(r *gin.RouterGroup, app *a.App) {
	g := r.Group("/repo", middlewars.Check_Role_Middleware(app))
	group := g.Group("/:repo_owner/:repo_name")

	group.GET("/meta", func(c *gin.Context) {
		ctx := c.Request.Context()
		repo_id := c.GetUint64(middlewars.Repo_id_feild)

		info, err := app.Git_db.Info(ctx, repo_id)
		if errors.Is(err, database.Error_Invalid) {
			c.JSON(bad_request(err))
			return
		} else if err != nil {
			c.JSON(internal_server_error())
			return
		}

		c.JSON(http.StatusOK, info)
	})

	group.GET("/info/refs", func(c *gin.Context) {
		ctx := c.Request.Context()
		service := c.Query("service")
		repo_id := c.GetUint64(middlewars.Repo_id_feild)

		if service != "git-receive-pack" && service != "git-upload-pack" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unsupported Service"})
			return
		}

		c.Header("Content-Type", fmt.Sprintf("application/x-%s-advertisement", service))
		c.Header("Cache-Control", "no-cache")

		err := git.Info_Refs(ctx, repo_id, service, c.Writer)
		if err != nil {
			c.JSON(internal_server_error())
			return
		}
	})

	group.POST("/git-upload-pack", func(c *gin.Context) {
		ctx := c.Request.Context()
		repo_id := c.GetUint64(middlewars.Repo_id_feild)

		c.Header("Content-Type", "application/x-git-upload-pack-result")
		c.Header("Cache-Control", "no-cache")

		err := git.Upload_Pack(ctx, repo_id, c.Request.Body, c.Writer)
		if err != nil {
			c.JSON(internal_server_error())
			return
		}
	})

	group.POST("/git-receive-pack", func(c *gin.Context) {
		ctx := c.Request.Context()
		repo_id := c.GetUint64(middlewars.Repo_id_feild)
		role := c.GetInt(middlewars.Role_feild)

		if role < roles.Owner {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Only owner can push"})
			return
		}

		err := git.Receive_Pack(ctx, repo_id, c.Request.Body, c.Writer)
		if errors.Is(err, database.Error_Invalid) {
			c.JSON(bad_request(err))
			return
		} else if err != nil {
			c.JSON(internal_server_error())
			return
		}
	})

	group.GET("/blob/*path", func(c *gin.Context) {
		repo_id := c.GetUint64(middlewars.Repo_id_feild)
		path := strings.Trim(c.Param("path"), "/")
		hash := c.Query("hash")

		blob, err := git.Get_Glob(repo_id, hash, path)
		if errors.Is(err, git.Error_Not_Found) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err == git.Error_Blob_Too_Large {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": err.Error()})
			return
		} else if err != nil {
			c.JSON(internal_server_error())
			return
		}

		blob_type := http.DetectContentType([]byte(blob))
		c.Header("Content-Type", blob_type)
		c.String(http.StatusOK, blob)
	})

	group.GET("/list/*path", func(c *gin.Context) {
		repo_id := c.GetUint64(middlewars.Repo_id_feild)
		path := strings.Trim(c.Param("path"), "/")
		hash := c.Query("hash")

		files, err := git.List_Dir(repo_id, hash, path)
		if errors.Is(err, git.Error_Not_Found) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err == git.Error_Path_Too_Deep {
			c.JSON(http.StatusRequestURITooLong, gin.H{"error": err.Error()})
			return
		} else if err != nil {
			c.JSON(internal_server_error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"files": files})
	})

	group.GET("/commits/:branch", func(c *gin.Context) {
		repo_id := c.GetUint64(middlewars.Repo_id_feild)
		branch := c.Param("branch")

		commits, err := git.Get_Commits(repo_id, branch)
		if errors.Is(err, git.Error_Not_Found) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err != nil {
			c.JSON(internal_server_error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"commits": commits})
	})

  group.GET("/branches", func(c *gin.Context) {
    repo_id := c.GetUint64(middlewars.Repo_id_feild)

    branches, err := git.Get_Branches(repo_id)
    if err != nil {
      c.JSON(internal_server_error())
      return
    }

    c.JSON(http.StatusOK, gin.H { "branches": branches })
  })
}
