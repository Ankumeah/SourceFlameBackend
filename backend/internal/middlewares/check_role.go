package middlewars

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

func Check_Role_Middleware(app *a.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		repo_name := c.Param("repo_name")
		repo_owner := c.Param("repo_owner")

		header := c.GetHeader("Authorization")
		if header == "" {
			c.Set("role", roles.Viewer)
		}

		var user_id uint64
		var pat_id uint64

		parts := strings.Split(header, " ")
		if len(parts) != 2 && header != "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Incorrct amount segments in Authorization header"})
			return
		} else if parts[0] != pat_auth_type && header != "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Wrong auth type, wanted: " + pat_auth_type})
			return
		} else if header != "" {
			username, pat, err := parse_basic_auth(parts[1])
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			user_id, err = app.User_db.Get_Id(ctx, username)
			if err == database.Error_invalid_user {
				c.AbortWithStatusJSON(bad_http_request(err))
				return
			} else if err != nil {
				c.AbortWithStatusJSON(internal_server_error())
				return
			}

			pat_id, err = app.PAT_db.Validate_PAT(ctx, user_id, pat)
			if errors.Is(err, database.Error_Invalid) {
				c.AbortWithStatusJSON(bad_http_request(err))
			} else if err != nil {
				c.AbortWithStatusJSON(internal_server_error())
				return
			}
		}

		owner_id, err := app.User_db.Get_Id(ctx, repo_owner)
		if err == database.Error_Invalid {
			c.AbortWithStatusJSON(bad_http_request(err))
			return
		}

		repo_id, err := app.Git_db.Get_Id(ctx, owner_id, repo_name)
		if err == database.Error_invalid_repo {
			c.AbortWithStatusJSON(bad_http_request(err))
			return
		} else if err != nil {
			c.AbortWithStatusJSON(internal_server_error())
			return
		}

		route_end := filepath.Base(c.Request.URL.Path)
		if route_end == "git-receive-pack" && user_id == 0 {
			c.Header("WWW-Authenticate", app.Settings.AUTH_CHALLENGE_HEADER)
			c.AbortWithStatusJSON(unauthorised_request())
			return
		}

		if user_id == owner_id {
			c.Set(Role_feild, roles.Owner)
		} else {
			c.Set(Role_feild, roles.Viewer)
		}

		app.PAT_db.Update_Use(ctx, pat_id)
		c.Set(User_id_feild, user_id)
		c.Set(Repo_id_feild, repo_id)
		c.Next()
	}
}
