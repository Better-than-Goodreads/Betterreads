package application

import (
	"github.com/betterreads/internal/domains/users/controller"
	"github.com/betterreads/internal/domains/users/repository"
	"github.com/betterreads/internal/domains/users/service"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
	port   string
}

func NewRouter(port string) *Router {
	r := gin.Default()
	addUsersHandlers(r)

	return &Router{
		engine: r,
		port:   port,
	}
}

func addUsersHandlers(r *gin.Engine) {
	rp := repository.NewMemoryDatabase()
	us := service.NewUsersService(rp)
	uc := controller.NewUsersController(us)

	r.POST("/users/register-first", uc.RegisterFirstStep)
    r.POST("/users/register-second", uc.RegisterSecondStep)
	r.POST("/users/login", uc.LogIn)
	r.GET("/users", uc.GetUsers)
}

func (r *Router) Run() {
	r.engine.Run(r.port)
}
