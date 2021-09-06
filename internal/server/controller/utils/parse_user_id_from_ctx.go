package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func ParseUserIDFromCtx(ctx *gin.Context) uint {
	sID := ctx.Param("user_id")
	uID, _ := strconv.ParseUint(sID, 10, 32)
	return uint(uID)
}
