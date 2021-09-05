package home

import (
	"github.com/Viverov/guideliner/internal/cradle"
	"github.com/gin-gonic/gin"
)

type Controller struct{}

func NewHomeController() *Controller {
	return &Controller{}
}

func (h *Controller) Init(router *gin.Engine, cradle *cradle.Cradle, prefix string) {
	router.GET(prefix + "/ping", pingHandler)
}

func pingHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
