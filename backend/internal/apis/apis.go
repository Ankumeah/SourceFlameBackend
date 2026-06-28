package apis

import (
	a "github.com/Ankumeah/SourceFlameBackend/internal/app"

	"github.com/gin-gonic/gin"
)

func Apis(r *gin.RouterGroup, app *a.App) {
	login(r, app)
	repos(r, app)
	repo(r, app)
	session(r, app)
	user(r, app)
	pat(r, app)
	ping(r)
}
