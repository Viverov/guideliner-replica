package server

import (
	"github.com/Viverov/guideliner/internal/cradle"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
	cradle *cradle.Cradle
}

func (s *Server) Run() {
	err := s.engine.Run()
	if err != nil {
		panic(err.Error())
	}
}

func Init(cradle *cradle.Cradle) *Server {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	return &Server{
		engine: r,
		cradle: cradle,
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
