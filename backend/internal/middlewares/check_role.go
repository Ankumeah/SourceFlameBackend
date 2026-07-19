package middlewares

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/database"
	"github.com/Ankumeah/SourceFlameBackend/internal/roles"

	"github.com/gin-gonic/gin"

	"errors"
	"net/http"
	"path/filepath"
	"strings"
)

func CheckRoleMiddleware(app *a.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		repoName := c.Param("repo_name")
		repoOwner := c.Param("repo_owner")

		header := c.GetHeader("Authorization")
		if header == "" {
			c.Set("role", roles.Viewer)
		}

		var userId uint64
		var patId uint64

		parts := strings.Split(header, " ")
		if len(parts) != 2 && header != "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Incorrect amount of segments in Authorization header"})
			return
		} else if parts[0] != patAuthType && header != "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Wrong auth type, wanted: " + patAuthType})
			return
		} else if header != "" {
			username, pat, err := parseBasicAuth(parts[1])
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			userId, err = app.UserDb.GetId(ctx, username)
			if errors.Is(err, database.ErrInvalidUser) {
				c.AbortWithStatusJSON(badHttpRequest(err))
				return
			} else if err != nil {
				c.AbortWithStatusJSON(internalServerError())
				return
			}

			patId, err = app.PATDb.ValidatePAT(ctx, userId, pat)
			if errors.Is(err, database.ErrInvalid) {
				c.AbortWithStatusJSON(badHttpRequest(err))
			} else if err != nil {
				c.AbortWithStatusJSON(internalServerError())
				return
			}
		}

		ownerId, err := app.UserDb.GetId(ctx, repoOwner)
		if errors.Is(err, database.ErrInvalid) {
			c.AbortWithStatusJSON(badHttpRequest(err))
			return
		}

		repoId, err := app.GitDb.GetId(ctx, ownerId, repoName)
		if errors.Is(err, database.ErrInvalidRepo) {
			c.AbortWithStatusJSON(badHttpRequest(err))
			return
		} else if err != nil {
			c.AbortWithStatusJSON(internalServerError())
			return
		}

		routeEnd := filepath.Base(c.Request.URL.Path)
		if routeEnd == "git-receive-pack" && userId == 0 {
			c.Header("WWW-Authenticate", app.Settings.AUTH_CHALLENGE_HEADER)
			c.AbortWithStatusJSON(unauthorisedRequest())
			return
		}

		if userId == ownerId {
			c.Set(RoleField, roles.Owner)
		} else {
			c.Set(RoleField, roles.Viewer)
		}

		app.PATDb.UpdateUse(ctx, patId)
		c.Set(UserIdField, userId)
		c.Set(RepoIdField, repoId)
		c.Next()
	}
}
