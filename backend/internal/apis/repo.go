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
	g := r.Group("/repo", middlewares.CheckRoleMiddleware(app))
	group := g.Group("/:repo_owner/:repo_name")

	group.GET("/meta", func(c *gin.Context) {
		ctx := c.Request.Context()
		repoId := c.GetUint64(middlewares.RepoIdField)

		info, err := app.GitDb.Info(ctx, repoId)
		if errors.Is(err, database.ErrInvalid) {
			c.JSON(badRequest(err))
			return
		} else if err != nil {
			c.JSON(internalServerError())
			return
		}

		c.JSON(http.StatusOK, info)
	})

	group.GET("/info/refs", func(c *gin.Context) {
		ctx := c.Request.Context()
		service := c.Query("service")
		repoId := c.GetUint64(middlewares.RepoIdField)

		if service != "git-receive-pack" && service != "git-upload-pack" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unsupported Service"})
			return
		}

		c.Header("Content-Type", fmt.Sprintf("application/x-%s-advertisement", service))
		c.Header("Cache-Control", "no-cache")

		err := git.InfoRefs(ctx, repoId, service, c.Writer)
		if err != nil {
			c.JSON(internalServerError())
			return
		}
	})

	group.POST("/git-upload-pack", func(c *gin.Context) {
		ctx := c.Request.Context()
		repoId := c.GetUint64(middlewares.RepoIdField)

		c.Header("Content-Type", "application/x-git-upload-pack-result")
		c.Header("Cache-Control", "no-cache")

		err := git.UploadPack(ctx, repoId, c.Request.Body, c.Writer)
		if err != nil {
			c.JSON(internalServerError())
			return
		}
	})

	group.POST("/git-receive-pack", func(c *gin.Context) {
		ctx := c.Request.Context()
		repoId := c.GetUint64(middlewares.RepoIdField)
		role := c.GetInt(middlewares.RoleField)

		if role < roles.Owner {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Only owner can push"})
			return
		}

		err := git.ReceivePack(ctx, repoId, c.Request.Body, c.Writer)
		if errors.Is(err, database.ErrInvalid) {
			c.JSON(badRequest(err))
			return
		} else if err != nil {
			c.JSON(internalServerError())
			return
		}
	})

	group.GET("/blob/*path", func(c *gin.Context) {
		repoId := c.GetUint64(middlewares.RepoIdField)
		path := strings.Trim(c.Param("path"), "/")
		hash := c.Query("hash")

		blob, err := git.GetBlob(repoId, hash, path)
		if errors.Is(err, git.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if errors.Is(err, git.ErrBlobTooLarge) {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": err.Error()})
			return
		} else if err != nil {
			c.JSON(internalServerError())
			return
		}

		blobType := http.DetectContentType([]byte(blob))
		c.Header("Content-Type", blobType)
		c.String(http.StatusOK, blob)
	})

	group.GET("/list/*path", func(c *gin.Context) {
		repoId := c.GetUint64(middlewares.RepoIdField)
		path := strings.Trim(c.Param("path"), "/")
		hash := c.Query("hash")

		files, err := git.ListDir(repoId, hash, path)
		if errors.Is(err, git.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if errors.Is(err, git.ErrPathTooDeep) {
			c.JSON(http.StatusRequestURITooLong, gin.H{"error": err.Error()})
			return
		} else if err != nil {
			c.JSON(internalServerError())
			return
		}

		c.JSON(http.StatusOK, gin.H{"files": files})
	})

	group.GET("/commits/:branch", func(c *gin.Context) {
		repoId := c.GetUint64(middlewares.RepoIdField)
		branch := c.Param("branch")

		commits, err := git.GetCommits(repoId, branch)
		if errors.Is(err, git.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err != nil {
			c.JSON(internalServerError())
			return
		}

		c.JSON(http.StatusOK, gin.H{"commits": commits})
	})

	group.GET("/branches", func(c *gin.Context) {
		repoId := c.GetUint64(middlewares.RepoIdField)

		branches, err := git.GetBranches(repoId)
		if err != nil {
			c.JSON(internalServerError())
			return
		}

		c.JSON(http.StatusOK, gin.H{"branches": branches})
	})

  group.GET("/blame/*path", func(c *gin.Context) {
    repoId := c.GetUint64(middlewares.RepoIdField)
		path := strings.Trim(c.Param("path"), "/")
		hash := c.Query("hash")

    blame, err := git.GetBlame(repoId, hash, path)
    if errors.Is(err, git.ErrInvalidCommitHash)  {
      c.JSON(badRequest(err))
      return
    } else if errors.Is(err, git.ErrBlobNotFound) {
      c.JSON(badRequest(err))
      return
    } else if err != nil {
      c.JSON(internalServerError())
      return
    }

    c.JSON(http.StatusOK, gin.H{"blame": blame})
  })
}
