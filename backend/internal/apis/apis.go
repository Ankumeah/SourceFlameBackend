package apis

import (
	"github.com/Ankumeah/DeltaBase/internal/session_store"

	"github.com/gin-gonic/gin"
)

func Apis(r *gin.RouterGroup, store *session_store.Session_store) {
  login(r, store)
}
