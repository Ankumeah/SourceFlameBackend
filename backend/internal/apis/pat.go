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
	group := r.Group("/pat", middlewares.SelfAuthMiddleware(*app.UserDb, false))

	group.POST("/:pat_name", func(c *gin.Context) {
		ctx := c.Request.Context()
		patName := c.Param("pat_name")
		userId := c.GetUint64(middlewares.UserIdField)

		pat, err := app.PATHandler.GeneratePAT()
		if err != nil {
			c.JSON(internalServerError())
			return
		}

		_, err = app.PATDb.AddPAT(ctx, userId, pat, patName)
		if errors.Is(err, database.ErrSafe) {
			c.JSON(badRequest(err))
			return
		} else if err != nil {
			c.JSON(internalServerError())
			return
		}

		c.JSON(http.StatusOK, gin.H{"PAT": pat})
	})

	group.DELETE("/:pat_name", func(c *gin.Context) {
		ctx := c.Request.Context()
		userId := c.GetUint64(middlewares.UserIdField)
		patName := c.Param("pat_name")

		patId, err := app.PATDb.GetId(ctx, userId, patName)
		if errors.Is(err, database.ErrInvalidPat) {
			c.JSON(badRequest(err))
			return
		} else if err != nil {
			c.JSON(internalServerError())
			return
		}

		if err = app.PATDb.DeletePAT(ctx, patId); err != nil {
			c.JSON(internalServerError())
			return
		}

		c.Status(http.StatusOK)
	})

	group.GET("/:pat_name", func(c *gin.Context) {
		ctx := c.Request.Context()
		userId := c.GetUint64(middlewares.UserIdField)
		patName := c.Param("pat_name")

		patId, err := app.PATDb.GetId(ctx, userId, patName)
		if errors.Is(err, database.ErrInvalidPat) {
			c.JSON(badRequest(err))
			return
		} else if err != nil {
			c.JSON(internalServerError())
			return
		}

		patInfo, err := app.PATDb.Info(ctx, patId)
		if err != nil {
			c.JSON(internalServerError())
			return
		}

		c.JSON(http.StatusOK, patInfo)
	})

	group.GET("/all", func(c *gin.Context) {
		ctx := c.Request.Context()
		userId := c.GetUint64(middlewares.UserIdField)

		pats, err := app.PATDb.GetPATs(ctx, userId)
		if errors.Is(err, database.ErrInvalid) {
			c.JSON(badRequest(err))
			return
		} else if err != nil {
			c.JSON(internalServerError())
			return
		}

		c.JSON(http.StatusOK, gin.H{"PATs": pats})
	})
}
