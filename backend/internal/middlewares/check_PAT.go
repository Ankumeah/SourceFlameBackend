package middlewars

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"
	"github.com/Ankumeah/SourceFlameBackend/internal/database"

	"github.com/gin-gonic/gin"

	"log"
	"errors"
	"net/http"
	"strings"
)

func Check_PAT_Middleware(app *a.App) gin.HandlerFunc {
  return func(c *gin.Context) {
    ctx := c.Request.Context()
    header := c.GetHeader("Authorization")
    if header == "" {
      c.Set("authed", false)
      c.Next()
      return
    }

    parts := strings.Split(header, " ")
    if len(parts) != 2 {
      c.JSON(http.StatusBadRequest, gin.H { "error": "Incorrct amount segments in Authorization header" })
      c.Abort()
      return
    } else if parts[0] != pat_auth_type {
      c.JSON(http.StatusBadRequest, gin.H { "error": "Wrong auth type, wanted: " + pat_auth_type })
      c.Abort()
      return
    }

    username, pat, err := parse_basic_auth(parts[1])
    if errors.Is(err, bad_request) {
      c.JSON(http.StatusBadRequest, gin.H { "error": err.Error() })
      c.Abort()
      return
    } else if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server error" })
      c.Abort()
      log.Printf("Error while parseing pat: %v\n", err.Error())
    }

    owner_id, err := app.User_db.Get_Id(ctx, username)
    if err == database.Error_invalid_user {
      c.JSON(http.StatusBadRequest, gin.H { "error": "Invalid user" })
      c.Abort()
      return
    } else if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server error" })
      c.Abort()
      return
    }

    _, err = app.PAT_db.Validate_PAT(ctx, owner_id, pat)
    if errors.Is(err, database.Error_Invalid) {
      c.JSON(http.StatusBadRequest, gin.H { "error": err.Error() })
      c.Abort()
    } else if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H { "error": "Internal server error" })
      c.Abort()
      return
    } else {
      c.Set("authed", true)
      c.Set("user_id", owner_id)
      c.Next()
      return
    }
  }
}


