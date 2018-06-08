package api

import (
	"LongTM/basic/api/auth"
	"LongTM/basic/system"
	"github.com/gin-gonic/gin"
)

func InitApi(root *gin.RouterGroup, tkWorker *system.VideoWorker) {
	auth.NewAuthenServer(root, "auth")
}
