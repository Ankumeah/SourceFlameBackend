package apis

import (
	"github.com/Ankumeah/SourceFlameBackend/internal/database"

	"github.com/gin-gonic/gin"

	"errors"
	"net/http"
)

func internalServerError() (int, map[string]any) {
	return http.StatusInternalServerError, map[string]any{"error": "Internal server error"}
}
func badRequest(err error) (int, map[string]any) {
	return http.StatusBadRequest, map[string]any{"error": err.Error()}
}

func getUserId(c *gin.Context, userDb *database.UserDb, username string) (uint64, bool) {
	ctx := c.Request.Context()

	userId, err := userDb.GetId(ctx, username)
	if errors.Is(err, database.ErrInvalidUser) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user"})
		return 0, false
	} else if err != nil {
		c.JSON(internalServerError())
		return 0, false
	} else {
		return userId, true
	}
}
