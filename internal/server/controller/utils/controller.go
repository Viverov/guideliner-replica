package utils

import (
	"github.com/Viverov/guideliner/internal/cradle"
	"github.com/gin-gonic/gin"
)

type Controller interface {
	// Init should initialize all handlers of this endpoints
	Init(router *gin.Engine, cradle *cradle.Cradle, prefix string)
}
