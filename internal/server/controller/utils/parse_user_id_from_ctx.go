package utils

import (
	"github.com/gin-gonic/gin"
)

func ParseUserIDFromCtx(ctx *gin.Context) uint {
	sID, exists := ctx.Get("user_id")
	if !exists {
		return 0
	}
	return sID.(uint)
}
