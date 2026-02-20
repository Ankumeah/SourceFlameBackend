package apis

import (
	"github.com/Ankumeah/DeltaBase/internal/session_store"

	"github.com/gin-gonic/gin"
)

func Apis(r *gin.RouterGroup, store *session_store.Session_store) {
  external_login(r, store)
}

func Session_Needed_Apis(r *gin.RouterGroup, store *session_store.Session_store) {
  logout(r, store)
  renew_session(r, store)
  renew_JWT(r)
}

func JWT_Needed_Apis(r *gin.RouterGroup) {
}
