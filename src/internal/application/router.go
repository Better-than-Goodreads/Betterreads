package application

import (
	"github.com/gin-gonic/gin"
	"github.com/betterreads/internal/domains/users/controller"
	"github.com/betterreads/internal/domains/users/service"
	"github.com/betterreads/internal/domains/users/repository"
)

type Router struct {
	engine *gin.Engine
	port  string
}

func NewRouter(port string) *Router {
	r := gin.Default()


	addUsersHandlers(r)

	return &Router{
		engine: r,
		port: port,
	}
}

func addUsersHandlers(r *gin.Engine) {
	rp := repository.NewMemoryDatabase()
	us := service.NewUsersService(rp)
	uc := controller.NewUsersController(us)

	r.POST("/users", uc.CreateUser)
	r.GET("/users", uc.GetUsers)
}



func (r *Router) Run() {
	r.engine.Run(r.port)
}