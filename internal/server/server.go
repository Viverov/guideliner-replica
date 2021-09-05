package server

import (
	"github.com/Viverov/guideliner/internal/cradle"
	"github.com/Viverov/guideliner/internal/server/controller/home"
	"github.com/Viverov/guideliner/internal/server/controller/user"
	"github.com/Viverov/guideliner/internal/server/controller/utils"
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

	var controllers []utils.Controller

	// Create utils
	httpResponder := utils.NewHttpResponder(cradle.GetConfig().Env)

	// Add controllers
	controllers = append(controllers, home.NewHomeController())
	controllers = append(controllers, user.NewUserController(httpResponder))

	// Init all controllers
	for _, c := range controllers {
		c.Init(r, cradle)
	}

	return &Server{
		engine: r,
		cradle: cradle,
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
