package apis

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/database"
	"github.com/Ankumeah/SourceFlameBackend/internal/middlewares"

	"github.com/gin-gonic/gin"

	"errors"
	"net/http"
)

func pat(r *gin.RouterGroup, app *a.App) {
	group := r.Group("/pat", middlewars.Self_Auth_Middleware(*app.User_db, false))

	group.POST("/:pat_name", func(c *gin.Context) {
		ctx := c.Request.Context()
		pat_name := c.Param("pat_name")
		user_id := c.GetUint64("user_id")

		pat, err := app.PAT_Handler.Genrate_PAT()
		if err != nil {
			c.JSON(internal_server_error())
			return
		}

		_, err = app.PAT_db.Add_PAT(ctx, user_id, pat, pat_name)
		if errors.Is(err, database.Safe_Error) {
			c.JSON(bad_request(err))
			return
		} else if err != nil {
			c.JSON(internal_server_error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"PAT": pat})
	})

	group.DELETE("/:pat_name", func(c *gin.Context) {
		ctx := c.Request.Context()
		user_id := c.GetUint64("user_id")
		pat_name := c.Param("pat_name")

		pat_id, err := app.PAT_db.Get_Id(ctx, user_id, pat_name)
		if err == database.Error_invalid_pat {
			c.JSON(bad_request(err))
			return
		} else if err != nil {
			c.JSON(internal_server_error())
			return
		}

		if err = app.PAT_db.Delete_PAT(ctx, pat_id); err != nil {
			c.JSON(internal_server_error())
			return
		}

		c.Status(http.StatusOK)
	})

	group.GET("/:pat_name", func(c *gin.Context) {
		ctx := c.Request.Context()
		user_id := c.GetUint64("user_id")
		pat_name := c.Param("pat_name")

		pat_id, err := app.PAT_db.Get_Id(ctx, user_id, pat_name)
		if err == database.Error_invalid_pat {
			c.JSON(bad_request(err))
			return
		} else if err != nil {
			c.JSON(internal_server_error())
			return
		}

		pat_info, err := app.PAT_db.Info(ctx, pat_id)
		if err != nil {
			c.JSON(internal_server_error())
			return
		}

		c.JSON(http.StatusOK, pat_info)
	})

	group.GET("/all", func(c *gin.Context) {
		ctx := c.Request.Context()
		user_id := c.GetUint64("user_id")

		pats, err := app.PAT_db.Get_PATs(ctx, user_id)
		if err == database.Error_Invalid {
			c.JSON(bad_request(err))
			return
		} else if err != nil {
			c.JSON(internal_server_error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"PATs": pats})
	})
}
